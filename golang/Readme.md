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

### Get environment properties
```
curl 'localhost:3100/get?key=temperature'
curl 'localhost:3100/get?key=humidity'
```

### Set environment properties
```
curl 'localhost:3100/set?key=temperature&value=15'
curl 'localhost:3100/set?key=humidity&value=50'
```
