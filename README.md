# Klokr
CLI practice project written in Golang. Combine a todo list with a time management aspect.

Planned Feature list(per august release):

1. To Do list -[x]
2. Manual appclocking -[x]
3. End of week report -[]
4. Tracking of time on each app -[]
5. Autolaunching of apps (may use cronjobs for this) -[]
6. Dashboard UI, graphs of some kind (maybe) -[]

8. How to use:
   klokr -h #list all commands
   klokr -a someactivity -f N #-a Adds new activity, -f sets frequency N of that activity, this is the reccomended way to add data
   klokr -c someactivity #decrements the frequency of selected activity by 1
   klokr -d someactivity #deletes selected activity
   klokr -p someactivity #prints the specific clockeditem struct, not very interesting right now
   Options:
     klokr -sde #toggles set delete when empty, when say you clock in an activity with value of 1, deletes the clocked object 
10. 
     Docker: 
        docker build -t klokr .
        docker run testklokr -it
     Golang:
       https://go.dev/doc/tutorial/compile-install
       Installation script will be forthcoming

     
