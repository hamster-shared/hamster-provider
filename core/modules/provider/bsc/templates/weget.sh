wget   $(curl -s https://api.github.com/repos/bnb-chain/bsc/releases/latest |grep browser_ |grep mainnet |cut -d\" -f4)
mkdir "config"
mv mainnet.zip config
cd config
unzip mainnet.zip
rm -rf mainnet.zip