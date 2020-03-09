# Back Pressure 背压/反压


作者：张铁蕾
链接：https://www.zhihu.com/question/49618581/answer/117107570

## 问题

数据流从上游生产者向下游消费者传输的过程中，上游生产者速度大于下游速度，导致下游的 Buffer 溢出。

解决这个问题就需要流量控制，对应的方法有：

* Backpressure 背压
* Throttling 节流
* Buffer & Window 打包
* Callstack blocking 阻塞


## Backpressure

消费者需要多少，生产者就生产多少。 这有点类似于TCP里的流量控制，接收方根据自己的接收窗口的情况来控制接收速率，并通过反向的ACK包来控制发送方的发送速率。这种方案只对于cold Observable有效。cold Observable是那些允许降低速率的发送源，比如两台机器传一个文件，速率可大可小，即使降低到每秒几个字节，只要时间足够长，还是能够完成的。相反的例子就是音视频直播，速率低于某个值整个功能就没法用了（这种类似于hot Observable）。


## Throttling
 说白了就是丢弃。消费不过来，就处理其中一部分，剩下的丢弃。至于处理哪些和丢弃哪些，就有不同的策略，也就是sample (or throttleLast)、throttleFirst、debounce (or throttleWithTimeout)这三种。还是举音视频直播的例子，在下游处理不过来的时候，就需要丢弃数据包。

## Buffer & Window

buffer和window基本一样，只是输出格式不太一样。它们是把上游多个小包裹打成大包裹，分发到下游。这样下游需要处理的包裹的个数就减少了。


## Callstack blocking 阻塞

 是一种特殊情况，阻塞住整个调用链（Callstack blocking）。之所以说这是一种特殊情况，是因为这种方式只适用于整个调用链都在一个线程上同步执行，这要求中间的各个operator都不能启动新的线程。在平常使用中这种应该是比较少见的，因为我们经常使用subscribeOn或observeOn来切换执行线程，而且有些复杂的operator本身也会内部启动新的线程来处理。另外 ，如果真的出现了完全同步的调用链，前面的方法仍然有可能适用的，只不过这种阻塞的方式更简单，不需要额外的支持。


 ## 一些解释


然后，从细的方面解释一下sample，throttleFirst，debounce。以及onBackpressureBuffer，onBackpressureDrop，onBackpressureBlock和ConnectableObservable（multicast）。

sample就是throttleLast，采样。类比一下音频采样，8kHz的音频就是每125微秒采一个值。sample可以配置成，比如每100毫秒采样一个值，但100毫秒内上游可能过来很多值，选那个值呢，就是选最后那个值。所以它也叫throttleLast。

throttleFirst跟sample类似，比如还是每100毫秒采样一个值，但选这100毫秒内的第一个值。

debounce，也叫throttleWithTimeout，名字里就包含一个例子。比如，一个网络程序维护一个TCP连接，不停地收发数据，但中间没数据可以收发的时候，就有间歇。这段间歇的时间，可以称为idle time。当idle time超过一个预设值的时候，就算超时了（timeout），这个时候可能就需要把连接断开了。实际上一些做server端的网络程序就是这么工作的。每收发一个数据包之后，启动一个计时器，等待idle time过去之后的超时，如果计时器到时之前，又有收发数据包的行为，那么计时器重置，等待一个新的idle time。当计时器到时了，就time out了，这个连接就可以关闭了。debounce的行为，跟这个非常类似，可以用它来找到连续的收发事件之后idle time超时后的timeout事件。

最后还有一个新的问题需要说明。Backpressure有些Observable是支持的，有些不支持。但它们可以通过operator来转化。

onBackpressureBuffer，onBackpressureDrop，onBackpressureBlock就可以把一个不支持Backpressure的Observable转成一个支持Backpressure的Observable（即支持request请求）。但转完之后的策略不太相同。

onBackpressureBuffer是不丢弃数据的处理方式。把上游收到的全部缓存下来，等下游来请求再发给下游。相当于一个水库。但上游太快，就会buffer溢出。

onBackpressureDrop就是当上游来数据的时候，看下游有没有需求，有需求就发给下游，否则上游来的数据就丢掉。

onBackpressureBlock也是看下游有没有需求，下游没有需求，不丢弃，但试图堵住上游的入口（能不能真堵得住还得看上游的情况了），自己并不缓存。

相反，有时候一些operator也能把一个支持Backpressure的Observable变成一个不支持Backpressure的Observable。比如，ConnectableObservable就是这样。它类似于把一条河的主干，在下游分成若干支流（但不太一样的是每条支流的水量都跟主干一样，是拷贝的）。那么很好理解，下游某个支流想对上游产生背压，是不太可能的，它阻止不了水流流向其它支流。
