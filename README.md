# go-vaccinate

go-vaccinate is a simulator for a hypothetical virus. It uses the <a href="github.com/gizak/termui/v3">termui</a> library for the console graphics. 

<img src="./demo.gif" />

## Building

To build go-vaccinate, you need go 1.14 and modules.

```sh
$ make clean build
```

## Running

If you build go-vaccinate locally, you have the option of configuring the runtime behavior. Run with no options to display usage info.

```sh
$ ./dist/simulator
Usage: simulator [--terminal|--console]

--console will run the simulation and just print the results.
--terminal will run the simulation and display the results using a plot and table
```

The console option will give you a brief menu of options.

```sh
$ ./dist/simulator  --console
Please select command
>
    load  Load configuration from ~/.vaccinate
    run   Run simulation
    quit  Quit
```

You must first *load* the configuration. If ~/.vaccinate is missing, a default one will be created.  You can then *run* the configuration. This can run quickly. Changing the configuration will change the runtime behavior.

```sh
$ ./dist/simulator  --console
Please select command
> load
Please select command
> run
COLUMN                       VALUE
Common name
People                       100
Visits                       100
Infection rate               10
Infected count               6
Number of  times infected    37
Number of times cured        31
Please select command
```

The docker container comes with a sample configuration and runs the terminal simulator. Once started it will run continuously until 'q' is pressed.

```sh
$ docker run --rm -ti vaccinate 
```

