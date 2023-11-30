# TCP server protected from DDOS by PoW

## Build image
docker build -t powserver .   

## Run
docker run -p 8080:8080 -e "SECRET=s3cr37" powserver

SECRET value is used in order to sign the challenges.

## Implementation

### Simplifications
The server is stand-alone. For simplicity of implementation only in-mem storage is used for challenges persistence.
The drawback is once server is restarted it won't remember that the solution nonce was already used against the challenge so it can be flooded with requests until validity period of challenge is reached.
For production use there should be Redis or DynamoDB used.

### Challenge implementation
The challenge implementation is inspired on PoW used by BitCoin so there's a nonce to be found by a client which satisfies the following condition 

hash(nonce + "." + challenge) starts with 00000  
