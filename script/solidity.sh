
# build contract get abi
forge build -C contract/test.sol --extra-output-files abi

# generate go code for abi
abigen --abi out/test.sol/Storage.abi.json  --pkg main --type Storage --out Storage.go


