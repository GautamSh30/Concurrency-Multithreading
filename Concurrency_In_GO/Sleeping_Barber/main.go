package main


import (
   "context"
   "fmt"
   "math/rand"
   "time"


   "github.com/fatih/color"
)


type ShopConfig struct {
   SeatingCapacity int
   ArrivalRate     int
   HairCutDuration time.Duration
   OpenDuration    time.Duration
}


func main() {


   // better RNG (testable + thread safe)
   r := rand.New(rand.NewSource(time.Now().UnixNano()))


   color.Red("The Sleeping Barber is about to wake up")
   color.Red("-----------------------------------------")


   config := ShopConfig{
       SeatingCapacity: 2,
       ArrivalRate:     100,
       HairCutDuration: 1000 * time.Millisecond,
       OpenDuration:    10 * time.Second,
   }


   clientChan := make(chan string, config.SeatingCapacity)
   doneChan := make(chan struct{})


   shop := BarberShop{
       ShopCapacity:    config.SeatingCapacity,
       HairCutDuration: config.HairCutDuration,
       ClientsChan:     clientChan,
       BarbersDoneChan: doneChan,
   }


   color.Green("The shop is open for the day")


   // add barbers
   shop.addBarber("Nai")
   shop.addBarber("Nai1")
   shop.addBarber("Nai2")
   shop.addBarber("Nai3")
   shop.addBarber("Nai4")
   shop.addBarber("Nai5")
   shop.addBarber("Nai6")


   ctx, cancel := context.WithCancel(context.Background())


   closed := make(chan struct{})


   // shop closing timer
   go func() {
       time.Sleep(config.OpenDuration)


       cancel()
       shop.closeShopForDay()


       close(closed)
   }()


   // generate clients
   go func() {


       i := 1


       for {


           randomMilliseconds := r.Intn(2 * config.ArrivalRate)


           select {


           case <-ctx.Done():
               return


           case <-time.After(time.Millisecond *
               time.Duration(randomMilliseconds)):


               shop.addClient(fmt.Sprintf("Client #%d", i))
               i++
           }
       }
   }()


   <-closed
}






package main


import (
   "time"


   "github.com/fatih/color"
)


type BarberShop struct {
   ShopCapacity    int
   HairCutDuration time.Duration
   NumberOfBarbers int


   BarbersDoneChan chan struct{}
   ClientsChan     chan string
}


// addBarber starts a worker goroutine
func (shop *BarberShop) addBarber(barber string) {


   shop.NumberOfBarbers++


   go func() {


       color.Yellow("%s starts working and waits for clients.", barber)


       for {


           // blocking receive = sleeping barber
           client, shopOpen := <-shop.ClientsChan


           if !shopOpen {


               shop.sendBarberHome(barber)
               return
           }


           color.Green("%s starts haircut for %s", barber, client)


           shop.cutHair(barber, client)
       }
   }()
}


func (shop *BarberShop) cutHair(barber string, client string) {


   color.Red("%s is cutting %s's hair.", barber, client)


   time.Sleep(shop.HairCutDuration)


   color.Green("%s finished cutting %s's hair.",
       barber,
       client,
   )
}


func (shop *BarberShop) sendBarberHome(barber string) {


   color.Cyan("%s is going home after working hard.", barber)


   shop.BarbersDoneChan <- struct{}{}
}


func (shop *BarberShop) closeShopForDay() {


   color.Cyan("Closing shop for the day.")


   // owner closes channel
   close(shop.ClientsChan)


   // wait for all barbers
   for i := 0; i < shop.NumberOfBarbers; i++ {


       <-shop.BarbersDoneChan
   }


   close(shop.BarbersDoneChan)


   color.Green(
       "Barber shop is now closed and everyone has gone home.",
   )
}


func (shop *BarberShop) addClient(client string) {


   color.Cyan("!!!!! %s arrived at shop.", client)


   select {


   case shop.ClientsChan <- client:


       color.Yellow(
           "%s takes a seat in the waiting room.",
           client,
       )


   default:


       color.Red(
           "Waiting room is full, so %s leaves.",
           client,
       )
   }
}
