library(ggplot2)

Tx1000CPU <- read.csv(file = "~/Downloads/1000 Tx CPU.csv", stringsAsFactors = TRUE)


stackedData <- cbind(Tx1000CPU[1:1], stack(Tx1000CPU[c(14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34)]))
stackedData$RelativeTime = c(0,diff(stackedData$Time))


ggplot(data = stackedData, aes(Time, values, group = ind, col = ind)) + 
  geom_line() + ylab('Host CPU Usage %') + labs(colour='Node') +
  ggtitle('1000 Transactions - CPU - Peers') + theme(plot.title = element_text(hjust = 0.5)) + theme(axis.text.x=element_blank(), axis.ticks.x=element_blank())



Tx1000CPU <- read.csv(file = "~/Downloads/10000 Tx Memory.csv", stringsAsFactors = TRUE)


stackedData <- cbind(Tx1000CPU[1:1], stack(Tx1000CPU[c(14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34)]))
stackedData$RelativeTime = c(0,diff(stackedData$Time))


ggplot(data = stackedData, aes(Time, values, group = ind, col = ind)) + 
  geom_line() + ylab('Host Memory Usage') + labs(colour='Node') +
  ggtitle('10,000 Transactions - Memory - Peers') + theme(plot.title = element_text(hjust = 0.5)) + theme(axis.text.x=element_blank(), axis.ticks.x=element_blank())

