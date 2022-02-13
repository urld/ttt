# ttt - time tracking terminal

[![Build Status](https://travis-ci.org/urld/ttt.svg?branch=master)](https://travis-ci.org/urld/ttt)
[![Go Report Card](https://goreportcard.com/badge/github.com/urld/ttt)](https://goreportcard.com/report/github.com/urld/ttt)
[![GoDoc](https://godoc.org/github.com/urld/ttt?status.svg)](https://godoc.org/github.com/urld/ttt)
[![GitHub release](https://img.shields.io/github/release/urld/ttt.svg)](https://github.com/urld/ttt/releases/latest)

`ttt` allows you to track your working hours, and generates simple reports.

## report format draft

```text
    week        date  worked   delta   saldo
          2021-11-23   7h42m
          2021-11-24   7h42m
          2021-11-25   
          2021-11-26   8h 0m  +  18m  +3h12m
          2021-11-28     45m  -6h57m
2021-w47               8h45m  -6h39m
          2021-11-29   8h15m
          2021-11-30   8h45m
          2021-12-01   8h15m
          2021-12-02   8h30m
          2021-12-03   8h 0m
2021-w48              41h45m

          2021-12-06   8h15m
2021-w49               8h15m
```
