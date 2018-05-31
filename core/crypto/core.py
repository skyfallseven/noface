'''
Skyfall Seven
---------------
Introduction to Cryptography
Project Assignment
---------------
'''

# To seed the lower half of the L1 matrix
from random import seed, randint
from numpy.linalg import inv # to inverse a
from numpy import matmul # for encryption
from numpy import array  # for NumPy Array to List conversion

def file2str(fileName):
    '''
    Returns file contents as a string
    '''
    with open(fileName, 'r') as f:
        x = f.readline().strip('\n')
    return x


def str2ascii(string):
	# converts string to ascii array
	return [ord(x) for x in string]

def f2inest(float_arr):
	# converts floats to ints in a nested array
	for i in range(len(float_arr)):
		for j in range(len(float_arr)):
			float_arr[i][j] = int(float_arr[i][j])
	return float_arr

def float2int(float_arr):
	# same as above, but for one 1D array
	for i in range(len(float_arr)):
		float_arr[i] = int(float_arr[i])
	return float_arr

def ascii2str(arr):
	# converts ascii array to string
	return ''.join(chr(int(x)) for x in arr)

def genSeed(arr):
	# generates seed based on array of ints
	s = 0
	for i in arr:
		s += i
	return s

def negMod(arr, q):
	# modulos every entry by prime (q)  to get rid of negatives
	for x in range(len(arr)):
		arr[x] = (arr[x] % q)

def initA(n, s):
	'''
	n: length of the message
	s: randomly place the 1/0 in the lower half 
		of the matrix. generated from the plaintext
	returns L1 and L2 matrix
	'''

	aMtrx = [[0 for x in range(n)]for y in range(n)]
	aInv = [[0 for x in range(n)]for y in range(n)]	
	# Now go down diagonals and write a 1
	x = 0	#The index to place the 1 at
	for i in range(n):
		for j in range(n):
			if (i == x and j == x):
				aMtrx[i][j] = 1
		x += 1

	# Randomly fill the lower half
	x = 0
	seed(s)
	for i in range(n):
		for j in range(n):
			if (j < x):
				aMtrx[i][j] = (randint(0,n) % 2)
		x += 1

	# now inverse it
	aInv = inv(aMtrx)
	aInv = aInv.tolist()
	aInv = f2inest(aInv)
	
	return aMtrx, aInv

def enc_F(q, n, p, t):
	'''
	Core encryption logic
	q: Prime number, in this case N_ASCII
	n: length of plaintext message
	p: plaintext as array of ASCII numbers
	t: password as array of ASCII numbers
	'''
	retval = None
	l = [0 for x in range(len(p))]
	if (len(t) % 2 == 0):
		retval = p
	else:
		retval = l
	r = 0
	# Please see project writeup "Central Map" section
	for i in range(len(t)):
		if (r == 0):
			l[0] = ((p[0] + t[i]) % q)
			l[1] = ((p[1] + p[0] * l[0]) % q)
			l[2] = ((p[2] + p[0] * l[1]) % q)
			for j in range(3,n):
				if (((j+1) % 4) == 3) or (((j+1) % 4) == 2):
					l[j] = ((p[j] + p[0] * l[j-2]) % q)
				else:
					l[j] = ((p[j] + p[j-2] * l[0]) % q)
			r = 1
		else:
			p[0] = ((l[0] + t[i]) % q)
			p[1] = ((l[1] - p[0] * l[0]) %q)
			p[2] = ((l[2] - p[0] * l[1]) %q)

			for j in range(3,n):
				if (((j+1) % 4) == 3) or (((j+1) % 4) == 2):
					p[j] = ((l[j] - p[0] * l[j-2]) % q)
				else:
					p[j] = ((l[j] - p[j-2] * l[0]) % q)
			r = 0
	return retval

def dec_F(q, n, c, t):
	'''
	Core decryption logic 
	See project writeup for details on algorithm
	q: Prime, in this case 127
	n: length of ciphertext
	c: ciphertext represented in ASCII codes
	t: password represented in ASCII codes	
	'''
	r = 0
	l = [0 for x in range(len(c))]
	p = [0 for x in range(len(c))]
	if (len(t)%2 == 0):
		p = c
	else:
		r = 1
		l = c

	for i in range(len(t)-1, -1, -1):
		if (r == 0):
			l[0] = ((p[0] - t[i]) % q)
			l[1] = ((p[1] + p[0] * l[0]) % q)
			l[2] = ((p[2] + p[0] * l[1]) % q)
			for j in range(3,n):
				if (((j+1) % 4) == 3) or (((j+1) % 4) == 2):
					l[j] = ((p[j] + p[0] * l[j-2]) % q)
				else:
					l[j] = ((p[j] + p[j-2] * l[0]) % q)
			r = 1
		else:
			p[0] = ((l[0] - t[i]) % q)
			p[1] = ((l[1] - p[0] * l[0]) % q)
			p[2] = ((l[2] - p[0] * l[1]) % q)
			for j in range(3,n):
				if (((j+1) % 4) == 3 or ((j+1) % 4) == 2):
					p[j] = ((l[j] - p[0] * l[j-2]) % q)
				else:
					p[j] = ((l[j] - p[j-2] * l[0]) % q)
			r = 0

	return p

