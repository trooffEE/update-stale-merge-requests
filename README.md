## `update-stale-merge-requests` is just handy script to update MR on target branch so you don't have to do it manually 🔧

### Installation: 
```shell
go install github.com/trooffEE/update-stale-merge-requests@latest
```

### Usage
```shell
update-stale-merge-requests # ease as that!
```

### This project requires 2 inputs to work, it will request them on first launch
1. GitLab Host (for instance for URL https://gitlab.test.test - host is `gitlab.test.test` so pass it when requested by script
2. Personal Access Token (**api rights only**) Read about it [here](https://docs.gitlab.com/user/profile/personal_access_tokens), usually it should be created here `https://<YOUR_GITLAB_HOST>/-/user_settings/personal_access_tokens`

These inputs are required to perform actions with your rights on hosts MRs, **they are stored locally on machine** (can check source code it fairly straight forward on purpose)

Config, after first initialisation, is placed in `~/.config/update-stale-config.yaml`.
If something went wrong (e.g. You passed incorrect token, or host) you can delete config by path above so that script on next launch will perform initialisation logic again
