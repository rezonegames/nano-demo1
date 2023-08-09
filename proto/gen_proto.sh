protoc --go_out=. *.proto

#
client_path=/Users/zyh/work/cocos-code/cocos_creator_framework
output_path=${client_path}/assets/Script/example/proto/
protoc  --plugin=/Users/zyh/.nvm/versions/node/v14.20.1/bin/protoc-gen-ts_proto \
--ts_proto_opt=esModuleInterop=true --ts_proto_opt=importSuffix=.js --ts_proto_opt=outputPartialMethods=false --ts_proto_opt=outputJsonMethods=false \
--ts_proto_out=${output_path} -I=. ./*.proto


# test
#output_path1=${client_path}/assets/Script/example/proto1/
#protoc --plugin=protoc-gen-ts=/Users/zyh/.nvm/versions/node/v14.20.1/bin/protoc-gen-ts \
#--ts_opt=target=web \
#--js_out=import_style=commonjs,binary:${output_path1} ./*.proto
