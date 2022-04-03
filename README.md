# tun
Quickly share files with friends using ngrok

## > Installation

To install, you will need a few things

 - An ngrok account
 - ngrok binary for your system
 - Running linux (as of v0.0.1a. You can try windows, but it is untested)

### > Windows

### > Linux

Put linux in your path with this command:
```bash
sudo cp $PWD/ngrok /usr/bin/ngrok
```

Authenticate with ngrok (check out their site for your token):
```bash
ngrok auth <token>
```

You're good to go! Check out the usage by running:
```bash
./tun
```