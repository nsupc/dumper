# Dumper
A simple utility for downloading nation and region dumps from [NationStates](https://www.nationstates.net).
## Arguments
 - user-agent (--user-agent / -u [USER]): NS nation or email address for API identification [Required]
 - nations (--nations / -n): download nation dump [Default: false]
 - regions (--regions / -r): download region dump [Default: false]
 - nation output folder (--out-dir-nations / -N [DIR]): the folder that will store the nation dump [Default: current directory]
 - region output folder (--out-dir-regions / -R [DIR]): the folder that will store the region dump [Default: current directory]
 - decompress (--decompress / -d): decompress the gzip archives to xml files [Default: false]
 - dry run (--dry-run / -D): test output file creation without downloading data from NationStates [Default: false]
