# gif.py

import argparse
import imageio
from os import listdir
from os.path import isfile, join


def define_arguments():
    parser = argparse.ArgumentParser()
    parser.add_argument('-f', '--delay', help="Frames per second", default=20)
    parser.add_argument('-n', '--filename', help="Filename for GIF", default="out")
    parser.add_argument('-p', '--pause', help="Pause (seconds) at end of gif loop", default=2.0)
    return parser.parse_known_args()[0]


def main():
    args = define_arguments()
    img_files = [f for f in listdir('./out/') if isfile(join('out', f))]
    duration = [float(1.0 / int(args.delay)) for i in img_files]
    duration[-1] = float(args.pause)
    images = [imageio.imread('./out/' + f) for f in img_files]
    imageio.mimsave('./out/' + args.filename + '.gif', images, duration=duration)

main()
