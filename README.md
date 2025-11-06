# Klokr
CLI practice project written in Golang. Combine a todo list with a time management aspect.
How to use:
   klokr -h #list all commands
   klokr -a someactivity -f N #-a Adds new activity, -f sets frequency N of that activity, this is the reccomended way to add data
   klokr -c someactivity #decrements the frequency of selected activity by 1
   klokr -d someactivity #deletes selected activity
   klokr -p someactivity #prints the specific clockeditem struct, not very interesting right now
   Options:
     klokr -sde #toggles set delete when empty, when say you clock in an activity with value of 1, deletes the clocked object 
 
Docker: 
        docker build -t klokr .
        docker run testklokr -it
Golang:
       https://go.dev/doc/tutorial/compile-install
       Installation script will be forthcoming


Planned Feature list(per November release):

1. To Do list - [x] - Now implemented as a json file
2. Manual appclocking - [x] - Can clock in activies via CLI, not proper endpoint yet.
3. Databse/SQL - [] - I realize this might be a bit extra for this kind of project, but it's okay for practices sake.
4. RAG - [] - Again might be a bit too much, but some type of data enriching might be useful.
5. End of week report - [] - Might need another file for this, and figure out rotations.
6. Tracking of time on each app - []
7. Autolaunching of apps (may use cronjobs for this) - []
8. Dashboard UI, graphs of some kind (maybe) - [] - Dashboard Partially implemented using a basic blazor template, frontend needs to be redone, possibly in another framework.
9. 
     
