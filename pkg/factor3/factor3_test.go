package factor3_test

import (
	"os"
	"sync"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drornir/factor3/pkg/example"
	"github.com/drornir/factor3/pkg/factor3"
	"github.com/drornir/factor3/tests"
)

var globalEnvMutex sync.Mutex

func init() { tests.Init() }

func TestExampleConfig(t *testing.T) {
	tFileSys := afero.NewMemMapFs()
	{
		ex1File, err := tFileSys.Create("ex1.yaml")
		require.NoError(t, err)
		_, err = ex1File.WriteString(`
version: v0
log:
  level: debug
  format: text
github:
  app:
    client_id: "put-it-here"
    pem_file: "/path/to/pem/file"
    installation_id: "1234567890"
`)
		require.NoError(t, err)
	}

	testCases := []struct {
		name      string
		filename  string
		envPrefix string
		env       map[string]string
		flags     []string

		checks func(t *testing.T, value example.Config)
	}{
		{
			name:     "test_basics",
			filename: "ex1.yaml",
			env: map[string]string{
				"TEST_BASICS_GITHUB_TOKEN": "my-secret-token",
			},
			flags: []string{"--log-level", "info", "--string", "fromflag"},
			checks: func(t *testing.T, value example.Config) {
				assert.Equal(t, "v0", value.Version)
				assert.Equal(t, "info", value.Log.Level)
				assert.Equal(t, "text", value.Log.Format)
				assert.Equal(t, "my-secret-token", string(value.Github.Token))
				assert.Equal(t, "/path/to/pem/file", string(value.Github.App.PemFile))
				assert.Equal(t, "fromflag", value.String)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// env part - using the real one until i can figure out out to inject viper a fake env
			globalEnvMutex.Lock()
			defer globalEnvMutex.Unlock()
			for k, v := range tc.env {
				prevVal, prevValExists := os.LookupEnv(k)
				os.Setenv(k, v)
				defer func(k string) {
					if prevValExists {
						os.Setenv(k, prevVal)
					} else {
						os.Unsetenv(k)
					}
				}(k)
			}

			viperInstance := viper.New()
			viperInstance.SetFs(tFileSys)
			viperInstance.AllowEmptyEnv(true) // maybe I should set it for everyone. it's a legacy feature to turn it off

			flagset := pflag.NewFlagSet(tc.name, pflag.ContinueOnError)
			err := factor3.InitializeViper(factor3.InitArgs{
				Viper:       viperInstance,
				ProgramName: tc.name,
				CfgFile:     tc.filename,
			})
			require.NoError(t, err, "factor3.InitializeViper() on %q", tc.filename)

			// this is where the magic happens
			var conf example.Config
			loader, err := factor3.Bind(&conf, viperInstance, flagset)
			require.NoError(t, err, "factor3.Bind()")
			require.NoError(t, flagset.Parse(tc.flags), "pflag.Parse returned an error on %v", tc.flags)

			err = loader.Load()
			require.NoError(t, err, "factor3.Load()")

			// loading was a success! now let's assert on the values
			tc.checks(t, conf)
		})
	}
}
