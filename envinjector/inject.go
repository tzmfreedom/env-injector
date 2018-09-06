package envinjector

import "os"

// InjectEnviron injects environment variables from AWS Parameter Store and/or SecretsManager
func InjectEnviron() {
	d := getDefaultDecorator()
	if path := os.Getenv("ENV_INJECTOR_META_CONFIG"); path != "" {
		injectEnvironViaMetaConfig(path)
	} else {
		trace("no meta config path specified, skipping injection via meta config")
	}

	if name := os.Getenv("ENV_INJECTOR_SECRET_NAME"); name != "" {
		injectEnvironSecretManager(name, d)
	} else {
		trace("no secret name specified, skipping injection by SecretsManager")
	}

	if path := os.Getenv("ENV_INJECTOR_PATH"); path != "" {
		injectEnvironByPath(path, d)
	} else {
		trace("no parameter path specified, skipping injection by path")
	}

	if prefix := os.Getenv("ENV_INJECTOR_PREFIX"); prefix != "" {
		injectEnvironByPrefix(prefix)
	} else {
		trace("no parameter prefix specified, skipping injection by prefix")
	}
}
