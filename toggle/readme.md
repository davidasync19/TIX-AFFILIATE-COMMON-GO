# Toggle

Toggle features is needed to make sure we have the flexibility for our release strategy.

Therefore, we need to have the abstraction standardization on how to implement the toggle features, with this standardization we expect able to:
- Maintain our code readability because no if-else branching inside our code.
- Easy to identify where is our toggle features is implemented, because all the code is wrap under ToggleExecutor interface
- When we want to remove our toggle capability experiment, we just can easy replace 1 line code of ToggleHelper.Run(spec) with the final code that we want.

## How To Implement

Assume we want to have a capability toggle features to do some calculation for given number, if this toggle capability is active then the calculation will use a formula result = value * 2 but if the toggle is inactive it will return the same value as given result = value.

`ToggleExecutor` implementation that implement above statement:
```
type SampleExecutor struct{}

/*
toggle.ToggleExecutor[T any, V any, W any]
T is the argument for OnToggleOn and OnToggleOff
V is the expected result of OnToggleOn and OnToggleOff
W is the argument for IsToggleOn
*/
func NewSampleExecutor() toggle.ToggleExecutor[int, int, interface{}] {
    return &SampleExecutor{}
}

/*
w can be ignored by using type interface{} and parsing `nil`.
Some operations can be conducted to suit your needs. For example determine the IsToggleOn by querying from databas or get the value from env variables.
*/
func (e *SampleExecutor) IsToggleOn(w interface{}) bool {
    // some operations
    return w == nil
}

func (e *SampleExecutor) OnToggleOn(t int) int {
    return t * 2
}

func (e *SampleExecutor) OnToggleOff(t int) int {
    return t
}
```
Afterwards, you just need to create ToggleHelper and used this inside the line that you want run this toggle experiment.
```
sampleExecutor := SampleExecutor{}
sampleHelper := toggle.ToggleHelper[int64, int64, interface{}]{
	Executor: &sampleExecutor,
}

// some processes
t := 2
sampleHelper.Run(t, nil)
```
Since we set t = 2 and w = nil the `IsToggleOn` will return true, `OnToggleOn` will be called, and we will get 4 as the result (t * 2).

For further reference: https://borobudur.atlassian.net/wiki/spaces/EN/pages/2594832707/Toggle+Features
