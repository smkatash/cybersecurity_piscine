# stockholm
#
#### Introduction to file manipulation by creating a harmless malware. 
#
> Note: `This project was completed for educational purposes only.  You should never use
this type of program for malicious purposes.`
#
#### Development
This project is completed with the aim of understanding how malware works, specifically ransomware.
A specific feature of this type of program is its ability to spread through networks of hundreds of computers. In our case, the program affects a small portion of local files. It’s all about understanding how a fairly simple program works in order to better protect oneself from it.

### Program

***stockholm*** recursively encrypts files in directory **infection** with **openssl** encryption library and stores secret key in the directory with the source code. Source code is written in Golang. Encrypted files are added **.ft** extension at the end of the original name. Once decrypted the file is returned to the original state. Program provides several options, including decryption given the correct key. 

##### Requirements
- Project is deployed within Docker container with Debian distribution (Docker file is provided). 
- Program acts only on directory **infection** placed in the root directory of HOME.
- For tests, make sure to include files whose extensions have been affected by **wannacry**. See [ref1][Technical Analysis of WannaCry].

##### Program Options
- "–help" or "-h" to display the help.
- "–version" or "-v" to show the version of the program.
- "–reverse" or "-r" followed by the key entered as an argument to reverse the infection.
- "–silent" or "-s", in which case the program will not produce any output.

##### Installation

Build and run the docker container

```sh
docker-compose up --build
```

Compile and run the project (in the Docker container)

```sh
cd stockholm
make
./stockholm 
```

Run with options

```sh
./stockholm // to encrypt
./stockholm -r secret.key // to decrypt
```


[//]: # (These are reference links used in the body of this note and get stripped out when the markdown processor does its job. There is no need to format nicely because it shouldn't be seen. Thanks SO - http://stackoverflow.com/questions/4823468/store-comments-in-markdown-syntax)

   [Technical Analysis of WannaCry]: <https://logrhythm.com/blog/a-technical-analysis-of-wannacry-ransomware/>
