# dalle-cli

dalle-cli is a little CLI for Mac OS X that enables you to generate and download images via OpenAI DALL E right from
your command line.

## Installation

ibooks-notes-exporter is available on OS X (both Intel and M-series processors).
It's distributed via a [homebrew](https://brew.sh/) package manager.

Run these commands in your terminal

```shell

> brew tap 57uff3r/mac-apps
> brew install 57uff3r/mac-apps/dalle-cli

```

## Usage

First of all, you have to get an API key for OpenAI API. Log in into your OpenAI account and
generate the API key [here](https://platform.openai.com/account/api-keys)

Then go to command line and create env variable with API key:

```shell
export DALLE_API_KEY=YOUR_API_KEY_GOES_HERE
```

And now you are ready to use dalle-cli!

```shell
dalle-cli generate --description "Nigh in cyberpunk city" --howmany 5
```
This command is going to generate 5 pictures using your textual description "Nigh in cyberpunk city".

Obviously, --description flag contains image description and --hownany contains a number of images that will be
generated. I've limited max number of images to 20 in order to prevent you from spending too much money on
generating too many images. All generated images will be generated in seprate folder.


## Feedback and contribution

Your feedback and pull requests are much appreciated.
Feel free to send your comments and thoughts to [me@akorchak.software](mailto:me@akorchak.software)


## Changelog
**0.0.1**

Initial release

### See also

Based on https://github.com/astralservices/go-dalle/