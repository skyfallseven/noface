'''
NoFace Development Build
Cryptography Core
TODO: Implement Chunking
'''
import core, getopt, sys

PRIME = 53813

def init(pt, pw):
	#pw = "@3sj5kl2!!!jfaZ9)"

	# Convert to ASCII codes
	pt = core.str2ascii(pt)
	pw = core.str2ascii(pw)

	# Generate matrices
	l1, l2 = core.initA(len(pt), core.genSeed(pw))
	return pt, pw, l1, l2

def encrypt(pt, pw, l1, l2):
	# Start basic scrambling
	x = core.matmul(pt, l1)
	x = x.tolist()
	core.negMod(x, PRIME)

	# Run encrytpion
	x = core.enc_F(PRIME, len(x), x, pw)

	# Scramble it again
	y = core.matmul(x, l2)
	y = y.tolist()
	core.negMod(y, PRIME)
	
	return core.ascii2str(y)

def decrypt(ct, pw, l1, l2):
	# Unscramble it 
	y = core.matmul(ct, l1)
	y = y.tolist()
	core.negMod(y, PRIME)

	# Decrypt
	y = core.dec_F(PRIME, len(y), y, pw)

	# Unscramble it again
	plain = core.matmul(y, l2)
	plain = plain.tolist()
	core.negMod(plain, PRIME)

	return core.ascii2str(plain)

def main():
	pt, pw, l1, l2 = init(sys.argv[2], sys.argv[3])
	if sys.argv[1] == "enc":
		print(encrypt(pt, pw, l1, l2))
	else:
		print(decrypt(pt, pw, l1, l2))

if __name__ == "__main__":
	main()
