type Adder interface {
	Add(x, y float64) float64
}

func (e Engine) Add(x, y float64) float64 {
	return x + y
}

func NewEngine() *Engine {
	return &Engine{}
}

type Calculator struct {
	Adder adder 
}

func NewCalculator(a Adder) *Calculator {
	return &Calculator{Adder : a}
}

func (c Calculator) PrintAdd(x, y float64) {
	fmt.Println("Result: ", c.Adder.Add(x, y))
}

func main() {
	engine := NewEngine()
	calc := NewCalculator(engine)
	cal.PrintAdd()
}