## Getting started

### Starting the simulation
```
make up
```

### Terminating the simulation
```
make down
```

### Cleaning up the simulation
```
make cleanup
```

### Set environment properties

#### Set temperature
```
curl 'localhost:3100/set?key=temperature&value=15'
```

#### Set humidity
```
curl 'localhost:3100/set?key=humidity&value=50'
```