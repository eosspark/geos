## DB FOR EOS
 use [storm][1]

## Simple Example

### Simple Struct 
``` go
type Zoo struct {
	Id     int `storm:"id,increment"` // id or ID 
	Name    string
	Animal  struct {
	    Number      int
		Carnivore   struct {
			Lion    int
			Tiger   int
		} `storm:"unique"`
	} `storm:"inline"`
}

``` 

### Test Code

*** Repeat deposit***
``` go
var zoo Zoo
zoo.Name = "zoo"
zoo.Animal.Number = 10
zoo.Animal.Carnivore.Lion = 100
zoo.Animal.Carnivore.Tiger = 100

err := db.Insert(&zoo)
if err != nil {
    //
}
err = db.Insert(&zoo)   // err == nil
var zoos []Zoo
err = db.All(&zoos)     // len(zoos) == 1

zoo.Id++
err = db.INsert(&zoo)   // already exists(unique)

```


### NOTE
- `Please read the test file carefully first`
- `Repeat the same structure, the database will not care, and will not return failure`

[1]: https://github.com/asdine/storm
