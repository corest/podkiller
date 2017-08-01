# Podkiller

Podkiller kills pods, tagged with special tags. Supports random schedules, black/white lists with namespaces.

## Custom configuration

See config.toml for examples. To replace configuration, mount your volume as `/etc/pod-killer` with `config.toml` file inside