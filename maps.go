package main

var Levels = []*Map{
	// Test Level
	NewMap(`
||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||
|
|
|
|
|
|
|
|
|
|
|
|
|                                         B
|                                         B
|                                         B    BBBBBBBBB
|                                         B    BBBBBBBBB
|                                         B    BBBBBBBBB
|                                         B    BBBBBBBBB
|                                         B    BBBBBBBBB
|                                         B    BBBBBBBBB
|  S             ######                        BBBBBBBBB
|               #######                        BBBBBBBBB
|              ########                        BBBBBBBBB
######################################################################################
`, true),
	NewMap(`
||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||

            #      #     ##
            #      #     ##
            #      #
            ########     ##
            #      #     ##
            #      #     ##
            #      #     ##















######################################################################################
`, false),
}

func init() {
	var last *Map
	for i, lvl := range Levels {
		lvl.ID = i
		if last != nil && last.Continues {
			lvl.OffsetX = last.OffsetX + last.Width
		}

		last = lvl
	}
}
