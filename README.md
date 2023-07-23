# Factor 3

Factor no. Three of the Twelve Factor App is "Config"

They suggest to store the configuration in env vars. I completely agree.

However, it's left to us to define how the env vars are loaded into our app.

This project is heavily inspired by 

- The legendary `github.com/spf13/viper`, and
- backstage.io's `app-config.yaml`

Also, it tries to make it easier to integrate with K8S Secrets and ConfigMaps.

The idea is to declaratively define a big struct, and pass that struct to 
this app, and it returns your filled config.