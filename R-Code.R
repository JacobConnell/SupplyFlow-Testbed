library(ggplot2)

data <- read.csv(file = "~/Desktop/Results/write/10000TxCPU.csv", stringsAsFactors = TRUE)

jumped <- data[,seq(1, ncol(data), 4)]

stackedData <- cbind(jumped[1:1], stack(jumped))



ggplot(data = stackedData, aes(ind, values, group = NAME, col = NAME)) + 
  geom_line() + ylab('Host CPU Usage %') + xlab('Seconds') + labs(colour='Node') +
  ggtitle('10000 Transactions - CPU - Peers') + theme(plot.title = element_text(hjust = 0.5)) 




data <- read.csv(file = "~/Desktop/Results/write/10000TxMem.csv", stringsAsFactors = TRUE)

jumped <- data[,seq(1, ncol(data), 4)]

stackedData <- cbind(jumped[1:1], stack(jumped))


ggplot(data = stackedData, aes(ind, values, group = NAME, col = NAME)) + 
  geom_line() + ylab('Host Mem Usage') + xlab('Seconds') + labs(colour='Node') +
  ggtitle('10000 Transactions - Memory - Peers') + theme(plot.title = element_text(hjust = 0.5)) 







data <- read.csv(file = "~/Desktop/Results/Caliper.csv", stringsAsFactors = TRUE)

jumped <- data[,seq(1, ncol(data), 1)]

stackedData <- cbind(jumped[1:1], stack(jumped))


ggplot(data = data) + 
  geom_line(aes(x=factor(TPS), y=Latency, group=1))  + ylab('Avg Latency (s)') + xlab('Input TPS') + labs(colour='Latency') +
  ggtitle('Latency by Input TPS') + theme(plot.title = element_text(hjust = 0.5)) 


ggplot(data = data) + 
  geom_line(aes(x=factor(TPS), y=Throughput, group=1))   + ylab('Throughput (TPS)') + xlab('Input TPS') + labs(colour='Throughput') +
  ggtitle('Throughput by Input TPS') + theme(plot.title = element_text(hjust = 0.5))



data <- read.csv(file = "~/Desktop/Results/Loadings.csv", stringsAsFactors = TRUE)

jumped2 <- data[,seq(1, ncol(data), 1)]

stackedData <- cbind(jumped2[1:1], stack(jumped2))


ggplot(data = data) + 
  geom_line(aes(x=factor(Tx.No), y=Throughput, group=1), colour="green") + geom_line(aes(x=factor(Tx.No), y=Latency, group=1), colour="red") + geom_line(aes(x=factor(Tx.No), y=TPS, group=1), , colour="orange")+ ylab('Avg Latency (s)') + xlab('Transaction Number') + labs(colour='ind') +
  ggtitle('Latency by Input TPS') + theme(plot.title = element_text(hjust = 0.5)) 

ggplot(data = LoadingsData) + 
  geom_line(aes(x=factor(Tx.No), y=values, group=ind, colour=ind))+ ylab('') + xlab('Transaction Number') + labs(colour='Metric') +
  ggtitle('High-load Transaction Metrics') + theme(plot.title = element_text(hjust = 0.5)) 

ggplot(data = data) + 
  geom_line(aes(x=factor(TPS), y=Throughput, group=1))   + ylab('Throughput (TPS)') + xlab('Input TPS') + labs(colour='Metric') +
  ggtitle('Throughput by Input TPS') + theme(plot.title = element_text(hjust = 0.5))





data <- read.csv(file = "~/Desktop/DBDownRData.csv", stringsAsFactors = TRUE)

jumped <- DBDownRData[,seq(1, ncol(DBDownRData), 4)]

stackedData <- cbind(jumped[1:1], stack(jumped))


ggplot(data = DBDownRData, aes(x=factor(State, level=c("Pre", "During", "Post")), group=1)) + geom_line(aes(y=Latency, col="Latency(S)")) + geom_line(aes(y=Threshold/3.5, col="Threshold(TPS)")) + scale_y_continuous(name="Latency(S)", sec.axis = sec_axis(~.*3.5, name="Threshold(TPS)")) +
xlab('Attack Stage') + labs(colour='Metric') + ggtitle('Attack Profromance Metrics') + theme(plot.title = element_text(hjust = 0.5))


ggplot(data = DBDownRData, aes(x=factor(State, level=c("Pre", "During", "Post")), y=Value, fill=tag)) + geom_bar(stat="identity", position="dodge") 

#+ geom_line(aes(y=Threshold/3.5, col="Threshold(TPS)")) + scale_y_continuous(name="Latency(S)", sec.axis = sec_axis(~.*3.5, name="Threshold(TPS)")) + xlab('Attack Stage') + labs(colour='Metric') + ggtitle('Attack Profromance Metrics') + theme(plot.title = element_text(hjust = 0.5))


ggplot(data = DBDownRData, aes(x=factor(State, level=c("Pre", "During", "Post")),  y=Value, fill=tag)) + geom_bar(stat="identity", position="dodge") + xlab('Attack Stage') + labs(colour='Metric') + ggtitle('Attack Performance Metrics') + theme(plot.title = element_text(hjust = 0.5)) +
  ylab('Rate')







#Batch Test Data Results


BatchSizeData <- read.csv("~/Desktop/BatchSizeData.csv")

ggplot(data = BatchSizeData, aes(x=Batch.Size))  + geom_bar(aes(y=Throughput/20, fill="Throughput (TPS)"), stat="identity") + geom_line(aes(y=Latency, colour="Latency (S)")) + scale_y_continuous(name="Latency(S)", sec.axis = sec_axis(~.*20, name="Throughput (TPS)")) +
labs(colour='Metric') + ggtitle('Batch Size Performance Metrics') + theme(plot.title = element_text(hjust = 0.5)) +  
  scale_fill_manual(name = NULL, values = c("Throughput (TPS)" = "light green")) +
  scale_color_manual(name = NULL, values = c("Latency (S)" = "red"))

