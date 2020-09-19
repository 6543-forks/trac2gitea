# trac2gitea

`trac2gitea` is a command-line utility for migrating [Trac](https://trac.edgewall.org/) projects to [Gitea](https://gitea.io/).

## Scope
At present the following Trac data is converted:
* Trac users mapped onto Gitea usernames (can be customised by providing an explicit mapping) 
* Trac components, priorities, resolutions, severities, types and versions to Gitea labels (can be customised by providing an explicit mapping)
* Trac milestones to Gitea milestones
* Trac tickets to Gitea issues
  * Trac ticket attachments to Gitea issue attachments
  * Trac ticket comments to Gitea issue comments with markdown text conversion
  * Trac ticket component, priority, resolution, severity, type and version changes to Gitea issue label changes
  * Trac ticket milestone changes to Gitea issue milestone changes
  * Trac ticket owner changes to Gitea issue assignee changes
  * Trac ticket "close" and "reopen" status changes to Gitea issue equivalents
  * Trac ticket summary changes to Gitea issue title changes
  * Trac ticket labels to Gitea issue labels
  * Trac ticket and comment owners to Gitea issue assignees
* Trac Wiki pages to files in the Gitea wiki repository
  * Markdown text conversion
  * Preservation of Trac wiki page history as separate wiki repository commits
* Trac to Gitea markdown conversions (copes with most cases but some Trac constructs may, possibly of necessity, not translate perfectly)
  * link anchors
  * block quotes
  * code blocks (single and multi-line)
  * definition lists
  * Trac bold, italic and underlines to markdown equivalents
  * headings
  * lists - bulletted, numbered, lettered and roman numbered
  * `[br]` paragraph breaks
  * tables (basic support)
  * Trac links:
    * images
    * `[[url|text]]` style
    * `[url text]` style
    * `http://...` and `https://...` links
    * `htdocs:...` (files are stored in a `htdocs` subdirectory of the Gitea wiki repository)
    * `CamelCase` inter-wiki links
    * `wiki:...` inter-wiki links
    * `attachment:...` current ticket or wiki page attachment references
    * `attachment:...:ticket:...` ticket attachment references
    * `attachment:...:wiki:...` wiki attachment references (files are stored in a `attachments/<pageName>` subdirectory of the Gitea wiki repository)
    * `ticket:...` ticket references
    * `comment:...` current ticket comment references
    * `comment:...:ticket:...` ticket comment references
    * `milestone:...` milestone references
    * `changeset:...` changeset references
    * `source:...` source file references

## Requirements ##
The utility requires access to both the Trac and Gitea filestore.
It retrieves data directly from the Trac database and writes into the Gitea database.
Access to the Gitea project wiki is via by checking out the wiki git repository.

The Gitea project must have been created prior to the migration as must the Gitea project wiki if a Trac wiki is to be converted (this can however just consist of an empty `Home.md` welcome page).

## Usage
```
Usage: trac2gitea [options] <trac-root> <gitea-root> <gitea-user> <gitea-repo> [<user-map>] [<label-map>]
Options:
      --db-only                   convert database only
      --generate-maps             generate default user/label mappings into provided map files (note: no conversion will be performed in this case)
      --no-wiki-push              do not push wiki on completion
      --overwrite                 overwrite existing data
      --verbose                   verbose output
      --wiki-convert-predefined   convert Trac predefined wiki pages - by default we skip these
      --wiki-dir string           directory into which to checkout (clone) wiki repository - defaults to cwd
      --wiki-only                 convert wiki only
      --wiki-token string         password/token for accessing wiki repository (ignored if wiki-url provided)
      --wiki-url string           URL of wiki repository - defaults to <server-root-url>/<gitea-user>/<gitea-repo>.wiki.git
```

* `<trac-root>` is the root of the Trac project filestore containing the Trac config file in subdirectory `conf/trac.ini`
* `<gitea-root>` is the root of the Gitea installation
* `<gitea-user>` is the owner of the Gitea project being migrated to
* `<gitea-repo>` is the Gitea repository (project) name being migrated to
* `<user-map>` is a file containing mappings from Trac users to Gitea user names - see below
* `<label-map>` is a file containing mappings from Trac items to Gitea labels - see below

### User Mappings
A file mapping from Trac users onto Gitea usernames can be provided via the `<user-map>` parameter.
This is a text file containing lines of the form:
```
<trac-user> = <gitea-username>
```

A default version of the mapping file can be generated by providing the `--generate-maps` flag.
This will write the default mapping into the user mapping file but not perform any actual data conversions.
The file can then be reviewed and edited as appropriate and the actual conversion process run by removing the `--generate-maps` flag.

If the `<user-map>` parameter is omitted, the conversion will proceed using the default mapping.

The default mapping maps a Trac user onto a Gitea user where the Gitea user has any of the following:
* the same user name
* the same "full" name
* the same email address

Where no mapping exists for a Trac user (the user map contains a line `<trac-user> =`):
* the Gitea repository owner provided on the command line will be used as the author of any issues or comments
* any Trac tickets assigned to the user will be converted left unassigned in Gitea
* the Trac user will be recorded as the "original author" of any Gitea issues

Where a mapping exists for a Trac user, the mapped Gitea user will be used in all relevant issues, comments etc.

### Label Mappings
A file mapping from Trac component, priority, resolution, severity, type and version names onto Gitea label names can be provided via the `<label-map>` parameter.
This is a text file containing lines of the form: 
```
<label-type>:<trac-item-name> = <gitea-label-name>
```
where `<label-type>` must be one of `component`, `priority`, `resolution`, `severity`, `type` and `version`.

As with user mappings, a default version of the mapping file can be generated by providing the `--generate-maps` flag.
This will write the default mapping into the label mapping file but not perform any actual data conversions.
Again, the file can then be reviewed and edited then actual conversion process run by removing the `--generate-maps` flag.
The mapping file can be edited so that `<gitea-label-name>` is unset for any Trac item (e.g. `resolution:fixed =`), in which case no label will be created for the Trac item.

The default mapping maps a Trac item name onto a Gitea label of the same name whether or not the Gitea label already exists.

If the `<label-map>` parameter is omitted, the conversion will proceed using the default mapping.

## Limitations
The current code is written for `sqlite` (for both the Trac and Gitea databases).

Hopefully, the SQL used by the converter is fairly generic so porting to a different database type should hopefully not be particularly difficult.

For anyone using a different database, the database connections are created in:
  * Trac: `accessor/trac/defaultAccessor.go`, func `CreateDefaultAccessor`
  * Gitea: `accessor/gitea/defaultAccessor.go`, func `CreateDefaultAccessor`

Having changed these, try running the converter and see if any SQL breaks.

All trac database accesses are in package `accessor.trac` and all Gitea database accesses are in package `accessor.gitea`.

## Building
From the root of the source tree run:
```
make
```
This will build the application as an executable `trac2gitea` (in the source tree root directory) and run the tests.

To build the application itself without running the tests, use:
```
make build
```

Missing dependencies can be fetched using:
```
make deps
```