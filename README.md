# configamajig

You know when you have 20 apps and they all are basically set up the same way, but not 100% the same? This is for that problem.

Configamajig will allow you to slurp in a bunch of yaml/json configs, then mix and match them and finally apply them to any type of file.

Real life use case:
Helm uses value.yaml files, but you have to deploy your app to multiple environments. Pretty simple, just keep a seperate copy for each env right? No multiply that by X apps. Suddenly the problem becomes much trickier to manage. You can use configamajig to instead keep 1 file for each app and template it out each env. This also makes big changes much easier! All your apps have to switch from one url to another for an api? Cool, change it once at a global level and everyapp gets it.

Configamajig will also let you remap files into other properties. Want to dump multiple files under the same reference variable? No problem, just use mappings.

You need to sometimes use one value vs another in an env depending on another defined value? Or maybe you want to switch a bunch to values at once? Great! You can use go templating logic and do that too.

Confused why something doesn't seem to be template like you expect? configamajig has full property tracing so you can see which file changed what super easy!

## Commands

### Replace:
Feed in an template file + a config and out comes a fully filled out file

```configamajig replace -c <CONFIG.json> -i <INPUT> -o <OUTPUT>```
