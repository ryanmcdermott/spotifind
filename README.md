# Spotifind


## What is it?
Spotifind is an app that uses Spotify's API to connect one artist to another by a chain of related artists.

You start with an origin artist, and you go through their related artists, and each related artists’ related artists, and so on and so forth, until you get a chain that links your origin artist to your destination artist by how they are related.


## API
Spotifind uses Spotify's API, so creating and registering an application is necessary. Go to [Spotify's developer site](https://developer.spotify.com/my-applications) to login with your Spotify account, and create your application. Call it whatever you want and don't worry about having a URI redirect or a website to link to. Just type the name of your application and its description and you'll get your client id and client secret. Save these, as we will need them shortly!


## Installation
1. `git clone https://github.com/ryanmcdermott/spotifind.git`
2. `cd spotifind`
3. `export SPOTIFY_CLIENT_ID="(YOUR_CLIENT_ID)"
4. `export SPOTIFY_CLIENT_SECRET="(YOUR_CLIENT_SECRET)"
5. `make vendor_get`
6. `make build`
7. The executable should be available ./bin/ and can be run by the example below.


## Example
`./bin/spotifind -o "Led Zeppelin" -d "Taylor Swift"`


## TODO
* Allow users to specify the number of concurrent workers to run
* Add a '--silent' option 
* Better error handling when the app exceeds Spotify's API request limit.


## Contributing
Pull requests are much appreciated and accepted.


## License
Spotifind is released under the [MIT License] (http://www.opensource.org/licenses/MIT)


## Credits
Nick Saika's [incredible article](http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html) on concurrent workers for Golang.
----