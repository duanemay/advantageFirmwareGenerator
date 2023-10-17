#!/usr/bin/env perl
use strict;
use warnings;

my $layerName = "BINDINGS_base";
my $insideBinding = 0;

open IN, "<" , $ARGV[0] or die "Can't open $ARGV[0] : $!";

while (my $line = <IN>) {

    if ($line =~ /layer_(\S) \{/ || $line =~ /(keypad) \{/ || $line =~ /(fn) \{/ || $line =~ /(mod) \{/) {
        $layerName = "BINDINGS_$1";
    }

    if ( $insideBinding == 0 && $line =~ /bindings = <$/ ) {
        $insideBinding = 1;
        print $line;
        print "        $layerName\n";
    } elsif ( $insideBinding == 1 && $line =~ />;$/ ) {
        $insideBinding = 0;
        print $line;
    } elsif ( $insideBinding == 1) {
        ## Do Nothing
    } else {
        print $line;
    }

}