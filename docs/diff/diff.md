# Diff

This repo is a fork from [zkevm-node](https://github.com/0xPolygonHermez.zkevm-node). The puirpose of the fork is to implement tha Validium consensus, enabling data availability to be posted outside of L1.

In order to document the code diff the [diff2html-cli](https://www.npmjs.com/package/diff2html-cli) tool is used. An html file is included in the repo [here](./diff.html). This file has been generated running the following command:

```bash
PATH_TO_ZKEVM_NODE_REPO="/change/this"
diff -ruN \
-I ".*github.com\/0x.*" \
-x "*mock*" -x ".git" \
-x ".github" \
-x ".gitignore" \
-x ".vscode" \
-x "ci" \
-x "environments" \
-x "*.md" \
-x "*.html" \
-x "*.html" \
-x "*.json" \
-x "*.toml" \
-x "*.abi" \
-x "*.bin" \
-x "*.pb.go" \
-x "smartcontracts" \
-x "go.sum" \
-x "mock*.go" \
-x "*venv*" \
-x "/dist/" \
-x "/test/e2e/keystore" \
-x "/test/vectors/src/**/*md" \
-x "/test/vectors/src/**/*js" \
-x "/test/vectors/src/**/*sol" \
-x "/test/vectors/src/**/*sh" \
-x "/test/vectors/src/package.json" \
-x "/test/contracts/bin/**/*.bin" \
-x "/test/contracts/bin/**/*.abi" \
-x "/tools/datastreamer/*.bin" \
-x "/test/datastreamer/*.db/*" \
-x "/test/*.bin" \
-x "/test/*.db/*" \
-x "**/.DS_Store" \
-x ".vscode" \
-x ".idea/" \
-x ".env" \
-x "out.dat" \
-x "cmd/__debug_bin" \
-x ".venv" \
-x "*metrics.txt" \
-x "coverage.out" \
-x "*datastream.db*" \
${PATH_TO_ZKEVM_NODE_REPO} . | \
diff2html -i stdin -s side -t "zkEVM node vs CDK validium node</br><h2>zkevm-node version: v0.5.0-RC4<h2/>" \
-F ./docs/diff/diff.html
```

Note that some files are excluded from the diff to omit changes that are not very relevant