import grpc from '@grpc/grpc-js';
import protoLoader from '@grpc/proto-loader';
import { searchUsersBySkill } from './service.js';

const PROTO_PATH = './proto/search.proto';

const packageDefinition = protoLoader.loadSync(PROTO_PATH);
const protoDescriptor = grpc.loadPackageDefinition(packageDefinition);
const searchProto = protoDescriptor.search;

const server = new grpc.Server();

server.addService(searchProto.SearchService.service, {
  SearchUsersBySkill: searchUsersBySkill
});

const PORT = 50051;
server.bindAsync(`0.0.0.0:${PORT}`, grpc.ServerCredentials.createInsecure(), () => {
  console.log(`SearchService running on port ${PORT}`);
  server.start();
});
