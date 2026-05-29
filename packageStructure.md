# Package Structure

### Main
Bare minimum behaviour here. Parsing CLI flags and coordinating other moving parts.

### Database
Actual SQL files and databases. Possibly some Go helper functions.

### Models
The Go representations of the necessary data structures. Should be entirely decoupled from front end rendering.

### Templates
All the html templates, css and javascript necessary for the web front end.

### Controller
Responsible for managing and routing http requests - either serving html or handling database changes and redirects

### Config
Project config and global settings

### App
The entrypoint that manages dependency injection and coordination within the whole project
