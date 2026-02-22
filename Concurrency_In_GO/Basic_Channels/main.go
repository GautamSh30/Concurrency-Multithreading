package main


import (
   "fmt"
   "math/rand"
   "time"


   "github.com/fatih/color"
)


const NumberOfPizzas = 10


var pizzasMade, pizzasFailed, total int


type Producer struct {
   data chan PizzaOrder
   quit chan chan error
}


type PizzaOrder struct {
   pizzaNumber int
   message     string
   success     bool
}


func (p *Producer) Close() error {
   ch := make(chan error)
   p.quit <- ch
   return <-ch
}


func makePizza(pizzaNumber int) *PizzaOrder {
   pizzaNumber++


   if pizzaNumber <= NumberOfPizzas {
       delay := rand.Intn(5) + 1
       color.Cyan("Received order #%d!", pizzaNumber)
       rnd := rand.Intn(12) + 1


       msg := ""
       success := false


       if rnd < 5 {
           pizzasFailed++
       } else {
           pizzasMade++
       }


       total++
       color.Yellow(
           "Making pizza #%d. It will take %d seconds...",
           pizzaNumber,
           delay,
       )
       time.Sleep(time.Duration(delay) * time.Second)


       if rnd <= 2 {
           msg = fmt.Sprintf(
               "!!! We ran out of ingredients for pizza #%d!",
               pizzaNumber,
           )
       } else if rnd <= 4 {
           msg = fmt.Sprintf(
               "!!! The cook quit while making pizza #%d!",
               pizzaNumber,
           )
       } else {
           success = true
           msg = fmt.Sprintf(
               "Pizza order #%d is ready!",
               pizzaNumber,
           )
       }


       p := PizzaOrder{
           pizzaNumber: pizzaNumber,
           message:     msg,
           success:     success,
       }
       return &p
   }


   return &PizzaOrder{
       pizzaNumber: pizzaNumber,
   }
}


func pizzaeria(pizzaMaker *Producer) {
   var i = 0
   for {
       currentPizza := makePizza(i)
       if currentPizza != nil {
           i = currentPizza.pizzaNumber


           select {
           case pizzaMaker.data <- *currentPizza:
           case quitChan := <-pizzaMaker.quit:
               close(pizzaMaker.data)
               // send nil error then close
               quitChan <- nil
               close(quitChan)
               return
           }
       }
   }
}


func main() {


   rand.Seed(time.Now().UnixNano())


   color.Green("Pizza Store is Open!")
   fmt.Println("---------------------")
   fmt.Println()


   pizzaJob := &Producer{
       data: make(chan PizzaOrder),
       quit: make(chan chan error),
   }


   go pizzaeria(pizzaJob)
   for i := range pizzaJob.data {
       if i.pizzaNumber <= NumberOfPizzas {
           fmt.Println(i.message)
           if i.success {
               color.Green(
                   "Order Number #%d is out for delivery!",
                   i.pizzaNumber,
               )
           } else {
               color.Red("Customer is really mad!")
           }
       } else {
           color.Magenta("Done making Pizzas!!")
           err := pizzaJob.Close()
           if err != nil {
               color.Red(
                   "!!! Error closing channel! %v",
                   err,
               )
           }
       }
   }


   fmt.Println("------------------------------")
   color.Green("Done for the Day!!!")
   fmt.Printf(
       "We made %d pizzas, failed %d pizzas, total attempts %d\n",
       pizzasMade,
       pizzasFailed,
       total,
   )
}
