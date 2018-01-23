# Testing Write Without Read in Client

Flood session response queue. See if the queue number keep increasing if client doesn't perform read on its side.

To test this property, we could use it like 2 users join chat, then 1 user keep sending messages, while another keep silent. Since the messages are broadcasted, if this property is true, then at some point server will output: "unable to relay message..." error. --> mungkin malah sebenarnya bakalan lebih bagus kalo si respQueue-nya unbuffered ya? --> jadi langsung ketahuan kalo ada delay --> yup2, patut dicoba, hehe...