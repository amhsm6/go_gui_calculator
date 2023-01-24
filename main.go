package main

import (
    "github.com/gotk3/gotk3/gtk"
    "log"
    "fmt"
    "math/big"
)

func main() {
    gtk.Init(nil)

    builder, err := gtk.BuilderNewFromFile("foo.glade")

    if err != nil {
        log.Panic(err)
    }

    obj, err := builder.GetObject("Number")

    if err != nil {
        log.Panic(err)
    }

    numberLabel := obj.(*gtk.Label)

    operations := map[rune]func(x, y big.Int) big.Int {
        '+': func(x, y big.Int) big.Int { return *x.Add(&x, &y) },
        '-': func(x, y big.Int) big.Int { return *x.Sub(&x, &y) },
        '*': func(x, y big.Int) big.Int { return *x.Mul(&x, &y) },
        '/': func(x, y big.Int) big.Int { if y.BitLen() == 0 { return y } else { return *x.Div(&x, &y) } },
        '=': func(x, y big.Int) big.Int { return y },
    }

    state := 0 // 0 - operation, 1 - number
    var lastRes big.Int
    lastOp := operations['=']
    currentSequence := "0"

    clear := func() {
        currentSequence = "0"
        numberLabel.SetLabel(currentSequence)

        state = 1
    }

    clearGlobal := func() {
        lastRes.SetUint64(0)
        numberLabel.SetLabel(lastRes.String())

        lastOp = operations['=']
        state = 0
    }

    changeSign := func() {
        if state == 0 {
            lastRes.Neg(&lastRes)
            numberLabel.SetLabel(lastRes.String())
        } else {
            if currentSequence[0] == '-' {
                currentSequence = currentSequence[1:]
            } else if currentSequence != "0" {
                currentSequence = "-" + currentSequence
            }

            numberLabel.SetLabel(currentSequence)
        }
    }

    signals := map[string]any {
        "clear": clear,
        "clearGlobal": clearGlobal,
        "changeSign": changeSign,
    }

    for i := 0; i < 10; i++ {
        currentNumber := i

        signals[fmt.Sprintf("number%d", i)] = func() {
            str_num := fmt.Sprintf("%d", currentNumber)

            if state == 0 {
                currentSequence = str_num
                numberLabel.SetLabel(currentSequence)

                state = 1
            } else {
                if currentSequence == "0" {
                    currentSequence = str_num
                } else {
                    currentSequence += str_num
                }

                numberLabel.SetLabel(currentSequence)
            }
        }
    }

    for operation, operationFunc := range(operations) {
        op := operation
        opFunc := operationFunc

        signals[string(op)] = func() {
            if state == 1 {
                var curSeqNum big.Int
                curSeqNum.SetString(currentSequence, 10)

                lastRes = lastOp(lastRes, curSeqNum)
                numberLabel.SetLabel(lastRes.String())

                state = 0
            }

            lastOp = opFunc
        }
    }

    builder.ConnectSignals(signals)

    winObj, _ := builder.GetObject("TopLevel")

    win := winObj.(*gtk.Window)

    win.SetDefaultSize(400, 400)

    win.Connect("destroy", gtk.MainQuit)

    win.ShowAll()

    gtk.Main()
}
