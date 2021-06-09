<!--
*** Thanks for checking out the Best-README-Template. If you have a suggestion
*** that would make this better, please fork the repo and create a pull request
*** or simply open an issue with the tag "enhancement".
*** Thanks again! Now go create something AMAZING! :D
***
***
***
*** To avoid retyping too much info. Do a search and replace for the following:
*** github_username, repo_name, twitter_handle, email, project_title, project_description
-->

<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->

<!-- PROJECT LOGO -->
<br />
<p align="center">
  <a href="https://dev.azure.com/c4ut/TimeClock/_git/auth-service">
    <img src="images/logo.png" alt="Logo" width="80" height="80">
  </a>

  <h3 align="center">TimeClock - Auth Service</h3>

  <p align="center">
    Microservice that provides an anti-corruption layer for authentication of microservices using keycloak
    <br />
    <a href="https://dev.azure.com/c4ut/TimeClock/_git/auth-service"><strong>Explore the docs »</strong></a>
    <!-- <br />
    <br />
    <a href="https://dev.azure.com/c4ut/TimeClock/_git/auth-service">View Demo</a>
    ·
    <a href="https://dev.azure.com/c4ut/TimeClock/_git/auth-service">Report Bug</a>
    ·
    <a href="https://dev.azure.com/c4ut/TimeClock/_git/auth-service">Request Feature</a>-->
  </p>
</p>

<!-- TABLE OF CONTENTS -->
<details open="open">
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <!-- <li><a href="#usage">Usage</a></li> -->
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <!-- <li><a href="#license">License</a></li> -->
    <li><a href="#contact">Contact</a></li>
    <!-- <li><a href="#acknowledgements">Acknowledgements</a></li> -->
  </ol>
</details>

<!-- ABOUT THE PROJECT -->
## About The Project

Auth service is a microservice that provides an anti-corruption layer for keycloak, it provides an API that makes it possible to establish a communication using REST and gRPC.

<!-- [![Product Name Screen Shot][product-screenshot]](https://example.com) -->
<!--
Here's a blank template to get started:
**To avoid retyping too much info. Do a search and replace with your text editor for the following:**
`github_username`, `repo_name`, `twitter_handle`, `email`, `project_title`, `project_description` -->

### Built With

- [Go Lang](https://golang.org/)
- List all: `go list -m all`

<!-- GETTING STARTED -->
## Getting Started

To get a local copy up and running follow these simple steps.

### Prerequisites

- Hiring a kubernetes cluster:
  - [AWS](https://aws.amazon.com/pt/eks/?whats-new-cards.sort-by=item.additionalFields.postDateTime&whats-new-cards.sort-order=desc&eks-blogs.sort-by=item.additionalFields.createdDate&eks-blogs.sort-order=desc)
  - [Azure](https://azure.microsoft.com/pt-br/services/kubernetes-service/)
  - [GCP](https://cloud.google.com/kubernetes-engine)

- [Kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)

- Create a secret for github docker registry

  ```sh
  kubectl create secret docker-registry regcred \
  --docker-server=$DOCKER_REGISTRY_SERVER \
  --docker-username=$DOCKER_USER \
  --docker-password=$DOCKER_PASSWORD \
  --docker-email=$DOCKER_EMAIL
  ```

- Create a secret with env credentials

  ```sh
  # file: credentials
  KEYCLOAK_REALM=keycloak_realm
  KEYCLOAK_CLIENT_ID=keycloak_client_id
  KEYCLOAK_CLIENT_SECRET=keycloak_client_secret
  KEYCLOAK_AUDIENCE=account
  ```

  `kubectl create secret generic auth-secret --from-env-file ./credentials`

### Deploy

- `kubectl apply -f ./k8s`

<!-- USAGE EXAMPLES -->
<!-- ## Usage

Use this space to show useful examples of how a project can be used. Additional screenshots, code examples and demos work well in this space. You may also link to more resources.

_For more examples, please refer to the [Documentation](https://example.com)_ -->

<!-- ROADMAP -->
## Roadmap

See the [open issues](https://dev.azure.com/c4ut/TimeClock/_backlogs/backlog/TimeClock%20Team/Epics) for a list of proposed features (and known issues).

<!-- CONTRIBUTING -->
## Contributing

Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

__Prerequisites__:

- Golang

  ```sh
  wget https://golang.org/dl/go1.16.2.linux-amd64.tar.gz
  rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.2.linux-amd64.tar.gz
  export PATH=$PATH:/usr/local/go/bin
  ```

- Docker and docker-compose

  ```sh
  sudo apt-get install docker docker-compose docker.io -y
  ```

- Environment

  ```sh
  # .env
  AUTH_GRPC_PORT=50051
  AUTH_REST_PORT=8088

  AUTHDB_DB=keycloak
  AUTHDB_USERNAME=keycloak
  AUTHDB_PASSWORD=password

  KEYCLOAK_USERNAME=admin
  KEYCLOAK_PASSWORD=Pa55w0rd
  KEYCLOAK_BASE_PATH=http://keycloak:8080
  KEYCLOAK_REALM=keycloak_realm
  KEYCLOAK_CLIENT_ID=keycloak_client_id
  KEYCLOAK_CLIENT_SECRET=keycloak_client_secret
  ```

__Installation__:

1. Clone the repo

   ```sh
   git clone https://dev.azure.com/c4ut/TimeClock/_git/auth-service.git
   ```

2. Run

   ```sh
   docker-compose up -d
   ```

3. Test

   ```sh
   go test -v -coverprofile cover.out ./...
   go tool cover -html=cover.out -o cover.html
   ```

__Installation in local kubernetes__:

1. Install [k3d](https://k3d.io/), [Kind](https://kind.sigs.k8s.io/) or similar
2. Install [Kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) and [Helm](https://helm.sh/)
3. Follow the steps of [Getting Started](#getting-started)
    - Connect to cluster
    - For the local keycloak, run:

      `helm install keycloak codecentric/keycloak`

    - For the local mongodb, run:

      `helm install mongo bitnami/mongodb`

      _add --set auth.enabled=false for no authentication_

    - [optional] For the local apm-server, run:
      `helm install apm-server elastic/apm-server`

    - finally, run:

      `kubectl apply -f k8s/`
<!-- LICENSE -->
<!-- ## License -->

<!-- Distributed under the MIT License. See `LICENSE` for more information. -->

<!-- CONTACT -->
## Contact

Coding4u - comercial@coding4u.com.br - [website](http://coding4u.com.br)

Project Link: [auth-service](https://dev.azure.com/c4ut/TimeClock/_git/auth-service)

<!-- ACKNOWLEDGEMENTS -->
<!-- ## Acknowledgements

* []()
* []()
* []() -->
