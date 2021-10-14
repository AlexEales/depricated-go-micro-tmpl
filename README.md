# Go Microservices Template

This is a template repository for kick-starting a Microservices project with Go, Docker,
Kubernetes, Helm, Hashicorp Vault, Protobuf (using buf), gRPC, gRPC-Gateway, Postgres, Redis and Kafka.

## Dependencies

Currently it is down to the user of the template to install the dependent libraries etc. required
for this project to function, however, below are links on how to install the various tools (up to
date at the time of writing this README):

- [Golang](https://golang.org/doc/install)
    - [Gomock](https://github.com/golang/mock)
- [Docker](https://docs.docker.com/desktop/)
- [Kubernetes](https://kubernetes.io/docs/tasks/tools/)
- [Helm](https://helm.sh/docs/intro/install/) - Used to simplify complicated deployments such as
  Postgres, Redis and Hashicorp Vault (can be customised with the `override-values.yaml` file).
- [Buf (tools for Protobuf, gRPC, gRPC-Gateway)](https://docs.buf.build/installation)

### Buf

Buf is a tool supported by Google (the authors of Protobuf) to allow for simple generation, linting
and documentation of you proto libraries. The main commands needed for running this projcted are as
follows:

- `buf generate` - This runs the buf generation scripts and produces the libraires specified in
  `buf.gen.yaml`, e.g. for gRPC a Go gRPC service file will be generated with the extension
  `<filename>_grpc.pb.go`.
- `buf lint` - Checks for linting errors on `.proto` file, very useful for consistency
- `buf push` - Pushes the buf library to the repository specified in `buf.yaml` in the `name` field
  (be sure to change this before using the project to whatever repo you want normally of the form: 
  `<buf username/org>/<repo_name>`).
- `buf mod update` - Unlikely to need this one but if you add any more deps into `buf.yaml` (which
  can be any public repo on the Buf repository) it should be added here to be included during
  compilation.

## TODO: Project Structure
