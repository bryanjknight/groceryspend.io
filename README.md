[![Browser Extension CICD](https://github.com/bryanjknight/groceryspend.io/actions/workflows/browser-extension.yml/badge.svg)](https://github.com/bryanjknight/groceryspend.io/actions/workflows/browser-extension.yml)
[![Server CICD](https://github.com/bryanjknight/groceryspend.io/actions/workflows/server.yml/badge.svg)](https://github.com/bryanjknight/groceryspend.io/actions/workflows/server.yml)
[![Web Portal CICD](https://github.com/bryanjknight/groceryspend.io/actions/workflows/web-portal.yml/badge.svg)](https://github.com/bryanjknight/groceryspend.io/actions/workflows/web-portal.yml)
[![Web WWW CICD](https://github.com/bryanjknight/groceryspend.io/actions/workflows/web-www.yml/badge.svg)](https://github.com/bryanjknight/groceryspend.io/actions/workflows/web-www.yml)

# What is GrocerySpend.io
Groceryspend.io is an app to allow a user to take receipts (either from an online shopping service like instacart) or from a image (from a traditional grocery store), and extract useful insights into how money is spent. For many families, groceries is the largest variable expensive month to month. The goal is help people better understand how they spend their money and make quality choices when balancing between health, environmental impact, and finances.


# Dev Setup
Below is a recommended list of settings and plugins to use in VSCode
## VSCode settings
```
{
    "workbench.activityBar.visible": true,
    "editor.tabSize": 2,

    "typescript.updateImportsOnFileMove.enabled": "always",
    "go.formatTool": "goimports",
    "go.lintTool": "golint",
    "[go]": {},
    "editor.defaultFormatter": null,
    "[javascript]": {
        "editor.defaultFormatter": "esbenp.prettier-vscode",
        "editor.formatOnSave": false
    },
    "[typescriptreact]": {
        "editor.defaultFormatter": "esbenp.prettier-vscode"
    },
    "[json]": {
        "editor.defaultFormatter": "esbenp.prettier-vscode"
    },
    "editor.minimap.enabled": false,
    "[jsonc]": {
        "editor.defaultFormatter": "esbenp.prettier-vscode"
    },
    "[typescript]": {
        "editor.defaultFormatter": "esbenp.prettier-vscode"
    },
    "redhat.telemetry.enabled": true
}
```

## VSCode plugins
- ms-azuretools.vscode-docker
- mikestead.dotenv
- dbaeumer.vscode-eslint
- eamodio.gitlens
- golang.go
- hashicorp.terraform
- ms-kubernetes-tools.vscode-kubernetes-tools
- ckolkman.vscode-postgres
- esbenp.prettier-vscode
- ms-python.python
- wayou.vscode-todo-highlight
- redhat.vscode-xml
- redhat.vscode-yaml

