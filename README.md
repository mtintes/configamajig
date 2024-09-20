# configamajig

You know when you have 20 apps and they all are basically configured the same way, but not 100% the same? This is for that problem.

Configamajig will allow you to slurp in a bunch of yaml/json configs, then mix and match them and finally apply them to any type of file.

Real life use case:
Helm uses value.yaml files, but you have to deploy your app to multiple environments. Pretty simple, just keep a seperate copy for each env right? Now multiply that by X number of apps. Suddenly the problem becomes much trickier to manage. You can instead keep 1 template file for each app and run configamajig on it. This also makes big changes much easier! All your apps have to switch from one url to another for any number of reasons? Cool, change it once at a global level and every app gets it.

Configamajig will also let you remap files into other properties. Want to dump multiple files under the same reference variable? No problem, just use mappings.

You need to sometimes use one value vs another in an env depending on another defined value? Or maybe you want to switch a bunch of values at once? Great! You can use go templating logic and [sprig functions](https://github.com/Masterminds/sprig) and do that too.

Confused why something doesn't seem to template like you expect? configamajig has full value tracing so you can see which file changed what value super easy!

## Commands

### Replace:
Feed in an template file + a config and out comes a fully filled out file

```configamajig replace -c <CONFIG.json> -i <INPUT> -o <OUTPUT>```
