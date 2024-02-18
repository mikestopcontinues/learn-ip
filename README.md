# What is an IP address?

An IP address (or Internet Protocol address) is a unique identifier assigned to a given machine. In the same way a physical address allows you to send mail to a particular place in the world, an IP address allows you to send requests to a particular device on a network.

## Public vs Private

Similar to physical addresses, IP addresses exist in both public and private networks. For example, you can mail a letter to the IRS in Kansas City, but you can't specify which office it goes to. Your letter will be routed in the mailroom to a particular agent.

In the same way, IP addresses come in two flavors:

- **Public IP addresses** are unique across the internet, and the devices they point to will typically be web servers or home modems.
- **Private IP addresses** are unique across a local network. Your laptop, cell phone, and game console will each have a unique local IP address, but none will have a public IP address.

## Static vs Dynamic

That said, IP addresses differ in a substantial way from physical addresses. IP addresses can change, and some do so frequently.

But which kinds of IP addresses change? And when?

- **Static IP addresses** are assigned to a particular device for the duration of a contract—sometimes lasting decades. Static IPs are typically used for web servers, as the reputation of a persistent IP address is valuable for security.
- **Dynamic IP addresses** are assigned by a DHCP server (Dynamic Host Configuration Protocol), and change often. Dynamic IPs are typically used for home networks, because it's cost-efficient for ISPs and more difficult for hackers to take advantage of.

## IPv4 vs IPv6

Originally, IP addresses were four bytes long, and were represented as a series of four numbers separated by dots. This was the IPv4 standard.

```text
192.168.0.1
```

And at the time, engineers thought that just over 4 billion addresses would be more than enough for all of the publicly-networked devices in the world.

But that's turned out not to be the case. In fact, we're rapidly running out of IP addresses. So in recent years, the world has been slowly migrating from the old format to a new format:

```text
2001:0db8:85a3:0000:0000:8a2e:0370:7334
```

The new format, IPv6, is sixteen bytes long, represented as eight groups of four hexadecimal digits each. (You may recognize hexadecimal from CSS hex colors codes like `#30fa99`.)

Because IPv6 addresses can be very long, you'll often find them shortened, to make them easier to read. Here's the rules:

- **Leading Zero Compressions**: You can drop the zeros at the start of every group.
- **Zero Compression**: Once per IP address, you can replace one or more zero group with `::`.

Some examples:

```bash
# leading zero compression
"0001:0002:0003:0004:0005:0006:0007:0008"
-->> "1:2:3:4:5:6:7:8"

# zero compression
"0000:1234:0000:0000:0000:0000:0000:5678"
-->> "0000:1234::5678"

# both compression rules
"0000:0000:0000:0000:0000:00B0:000B:001E"
-->> "::B0:0B:1E" 🫣
```

Whereas IPv4 supports 2^32 addresses, IPv6 supports 2^128 addresses. This is enough to assign a unique address to every device that exists today, as well as for a lifetime ahead. This allows for exciting new patterns for decentralization and "the internet of things."

Thus, we can categorize IPs one final way:

- [**IPv4 addresses**](https://en.wikipedia.org/wiki/Internet_Protocol_version_4) are still widely used, but they are limited in number, and will be mostly retired in the next decade.
- [**IPv6 addresses**](https://en.wikipedia.org/wiki/IPv6) are still an emerging standard, but they provide ample space to support humanity's very digital future.

## Special IP ranges

Both IPv4 and IPv6 have IP address ranges for dedicated tasks. While there's no need to memorize these addresses, they're worth exploring, as you'll encounter them again and again in your web development career.

### Loopback addresses

- IPv4: `127.*.*.*`
- IPv6: `::1`

Loopback IP addresses allow a device to communicate with itself. You probably think of this as `localhost`, but `localhost` is actually just shorthand for `127.0.0.1`. It's rare to use any other loopback address.

In IPv4, every address starting with `127.` is a loopback address. However IPv6 only supports one, `::1`. By the time they designed IPv6, they found that nobody really needed 17 million IP addresses for their local machine!

### Private network addresses

- IPv4: `192.168.*.*`, `10.*.*.*`, `172.16.*.*`
- IPv6: `fd*::`

In IPv4, any address starting with `192.168.` is private. You can reach such addresses within your local network, but no external traffic can reach them. Typically, your router will assign a IP matching `192.168.1.*`to each device on your network.

Many large organizations began to hit limits within this range. Not typically because they had more than 65 thousand devices (`256*256`) on their local network. But because you can link local networks together. And when you do, there are frequently conflicts.

To solve this, IPv6 replaced private network addresses with ULAs (Unique Local Addresses). ULAs split the IP address into a random ID for the organization and a random ID for the machine, ensuring that it's incredibly unlikely to have conflicts.

---

## Exercise: Hack Your Network

We've talked a lot about IP addresses. Let's see if you can find yours.

In this project, you're going to scan your local network for devices, then you'll also scan common ports to see which machines have servers running.

You'll be using some new Go packages, and some oldies-but-goodies:

- [fmt](https://pkg.go.dev/fmt)
- [net](https://pkg.go.dev/net)
- [os/exec](https://pkg.go.dev/os/exec)
- [sync](https://pkg.go.dev/sync)
- [time](https://pkg.go.dev/time)

### Step 0: What's a ping?

A network is a lot like the deep ocean. You can't see anything, and so if you need to know what's out there, you have to send out a signal and see what comes back.

The most basic signal is a "ping." This is a small packet of data that you send to a machine, and if it's on, it will send a response back.

To start, jump into your CLI and run the following command:

```bash
ping -c 1 -W 1 127.0.0.1
```

What do you see?

Unless something's gone completely wrong, you should see a response like this:

```text
PING 127.0.0.1 (127.0.0.1): 56 data bytes
64 bytes from 127.0.0.1: icmp_seq=0 ttl=64 time=0.075 ms

--- 127.0.0.1 ping statistics ---
1 packets transmitted, 1 packets received, 0.0% packet loss
round-trip min/avg/max/stddev = 0.075/0.075/0.075/0.000 ms
```

This is your computer sending a ping to itself, and receiving a response back. 0.0% packet loss means that the ping was successful.

Note, in the above command, we specify two arguments:

- `-c 1` tells the `ping` command to send only one packet.
- `-W 1` tells the `ping` command to wait only one second for a response.

If you omit these arguments, `ping` will run indefinitely, even if it's not getting a response. It can be useful for debugging, but it's not what we want for now.

## Step 1: Scan your local network

Now that you know what a ping is, let's write a Go program to send pings to every _potential_ device on your local network.

Write code that does the following:

1. Loop through every IP address matching `192.168.1.*`. (How many potential IPs are in that range?)
2. For each IP address, increment your `WaitGroup` and start a goroutine that defer's a `wg.Done` call.
3. Within each goroutine, send a ping to the IP address. You'll need the `exec` package for this.
4. Print every IP address that responds to your ping, and skip every IP address that doesn't.

If everything goes well, you should see a list of IP addresses that are currently active on your local network. Mine looks like this:

```text
192.168.1.1
192.168.1.12
192.168.1.251
```

## Step 2: Scan common ports

Now that you know which devices are active, let's see which ones are running servers. For this, you'll have to fill in `scanCommonPorts()`, then use it within your `main()` loop.

Here's the plan:

1. Configure a new `WaitGroup` for your `scanCommonPorts()` function.
2. Loop through every common port in our `commonPorts` map.
3. For each port, increment your `WaitGroup` and start a goroutine that defer's a `wg.Done` call.
4. Within each goroutine, use `net.DialTimeout()` to send a "tcp" request to your IP + port.
    - Use the `time` package to specify a timeout of 1 second.
    - Make sure you close your connection when your done with it.
5. Print every IP + port that responds to your ping, and skip every IP address that doesn't.
    - Include the description of the server on that port. For example "(FTP)" or "(VNC)".

If everything goes well, you should see a list of IP addresses and ports that are currently active on your local network. Mine looks like this:

```text
192.168.1.1
192.168.1.1:53 (DNS)
192.168.1.1:80 (HTTP)
192.168.1.251
192.168.1.251:80 (HTTP)
192.168.1.251:53 (DNS)
192.168.1.251:443 (HTTPS)
```

## Step 3: Hack your family's machines

Now that you know what everyone in your household is up to, it's time to start messing with them. Nothing more fun that DDoSing your kid's Minecraft server, right?

(Just kidding. Don't do that.)
