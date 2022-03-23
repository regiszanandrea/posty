# Definition
posty is very similar to Twitter, but it has far fewer features.

# Technologies

- Language: Golang
- Frameworks: Fiber and FX
- Database: MongoDB

# Setup
There are two ways to setup the project, the hard way will be more detailed on [setup](./docs/setup.md) file, on backstage docs.

The easy way will be describe in the next follow steps:

1. Install docker-compose: https://docs.docker.com/compose/install/#install-compose-on-linux-systems
2. Run the follow command and edit the environment variables as you wish:
```shell
cp configs/app/example.yaml configs/app/local.yaml
```
3. Run the follow command on the root project folder:
```shell
docker-compose up -d
```
4. Access http://127.0.0.1:3030

# Seed
To create some data on database, you can just run `make seed` that it will generate some fake users, posts and followers

# API
There are two files to use if you would like to see all endpoints from the application:
1. With file `openapi.json` you can copy and put at https://editor.swagger.io to see all endpoints on swagger ui
2. Import the file `posty.postman_collection.json` in your Postman application

# Developing

If you're using the docker-compose setup, it auto reloads your application automatically on saving any file.

# Testing

1. Install all applications dependencies
2. Run the follow command and edit the environment variables as you wish:
```shell
cp configs/app/example.yaml configs/app/testing.yaml
```
3. Run `make test` to execute all tests

# Contribute

- Every change should generate a Pull request to main branch

# Planning

Questions to Product Manager:
- Does quoted-posts are also be returned on this new page? 
- Does this should increment posts count?
- Should this be considered a post for the 5 posts per day restriction? 
- Should exist a restriction of the number of these posts per day? 
- Does a QA Engineer should be involved? to make a "Three Amigos" meeting, with his help
we could define better test scenarios 

Changes on the application:
- Change the endpoint to get user's posts that don't have `parent_id`
- Add a new repository and service method, and an endpoint to get all posts from user
- Should be added an index to an `content` field on posts collection, for a full-text search, as "@" is a reply. 

# Critique

Things that could be improved:
- Increase service layer test coverage, mainly with negative cases
- Make seeder accepts parameters
- Separate interfaces in smaller interfaces to make mocking easier
- Make a health check endpoint
- About scaling, I think MongoDB could scale very well, you can have some replicas to scale horizontally, also Golang it's very fast. Talking about infrastructure, could have
  a layer of cache on feed endpoint to avoid too much load on the database, also some parts could be done on Event-Sourcing architecture like when a user follows someone, 
  it could dispatch an event of the following and the user's domain could listen to this event to increment the follower's count, instead of doing in the same request as it is like now.
