public final class SwitchWide {
    static int dense(int key) {
        switch (key) {
            case 0: return 11;
            case 1: return 21;
            case 2: return 31;
            case 3: return 41;
            case 4: return 51;
            case 5: return 61;
            case 6: return 71;
            case 7: return 81;
            case 8: return 91;
            case 9: return 101;
            default: return -1;
        }
    }

    static int offsetDense(int key) {
        switch (key) {
            case 10: return 201;
            case 11: return 202;
            case 12: return 203;
            case 13: return 204;
            case 14: return 205;
            case 15: return 206;
            case 16: return 207;
            case 17: return 208;
            case 18: return 209;
            case 19: return 210;
            default: return -10;
        }
    }

    static int sparse(int key) {
        switch (key) {
            case -1000: return 5;
            case -7: return 6;
            case 0: return 7;
            case 3: return 8;
            case 100: return 9;
            case 10000: return 10;
            case 1000000: return 11;
            default: return -5;
        }
    }

    static int nested(int a, int b) {
        int r = 0;
        switch (a) {
            case 0:
                switch (b) {
                    case 0: r = 13; break;
                    case 1: r = 14; break;
                    case 2: r = 15; break;
                    default: r = 16; break;
                }
                break;
            case 1:
                switch (b) {
                    case 5: r = 17; break;
                    case 6: r = 18; break;
                    default: r = 19; break;
                }
                break;
            case 2:
                r = 20;
                break;
            default:
                r = 21;
                break;
        }
        return r;
    }

    static int loopSwitch(int limit) {
        int acc = 0;
        for (int i = 0; i < limit; i++) {
            switch (i) {
                case 0: acc += 1; break;
                case 1: acc += 2; break;
                case 2: acc += 3; break;
                case 3: acc += 4; break;
                case 4: acc += 5; break;
                case 5: acc += 6; break;
                default: acc += 7; break;
            }
            if (acc > 15) {
                break;
            }
        }
        int j = limit;
        while (j > 0) {
            acc += j;
            j -= 1;
        }
        return acc;
    }

    static int charSwitch(char c) {
        switch (c) {
            case 'a': return 31;
            case 'b': return 32;
            case 'z': return 33;
            default: return 34;
        }
    }

    static int byteSwitch(byte b) {
        switch (b) {
            case -2: return 41;
            case 0: return 42;
            case 5: return 43;
            default: return 44;
        }
    }

    static int shortSwitch(short s) {
        switch (s) {
            case -300: return 51;
            case 0: return 52;
            case 300: return 53;
            default: return 54;
        }
    }

    @SuppressWarnings("fallthrough")
    static int fallthrough(int key) {
        int r = 0;
        switch (key) {
            case 5:
                r += 1;
            default:
                r += 2;
            case 7:
                r += 4;
                break;
            case 9:
                r += 8;
                break;
        }
        return r;
    }

    static int noDefault(int key) {
        int r = 100;
        switch (key) {
            case 0:
                r = 101;
                break;
            case 1:
                r = 102;
                break;
        }
        return r;
    }

    static int deepIf(int value) {
        int r = 0;
        if (value > 100) {
            if (value > 1000) {
                if (value > 10000) {
                    r = 4;
                } else {
                    r = 3;
                }
            } else {
                if (value > 500) {
                    r = 2;
                } else {
                    r = 1;
                }
            }
        } else {
            if (value > 10) {
                if (value > 50) {
                    r = 8;
                } else {
                    r = 7;
                }
            } else {
                if (value > 0) {
                    r = 6;
                } else if (value == 0) {
                    r = 5;
                } else {
                    r = 9;
                }
            }
        }
        return r;
    }

    static int wideInts(int seed) {
        int v0 = 0;
        int v1 = 1;
        int v2 = 2;
        int v3 = 3;
        int v4 = 4;
        int v5 = 5;
        int v6 = 6;
        int v7 = 7;
        int v8 = 8;
        int v9 = 9;
        int v10 = 10;
        int v11 = 11;
        int v12 = 12;
        int v13 = 13;
        int v14 = 14;
        int v15 = 15;
        int v16 = 16;
        int v17 = 17;
        int v18 = 18;
        int v19 = 19;
        int v20 = 20;
        int v21 = 21;
        int v22 = 22;
        int v23 = 23;
        int v24 = 24;
        int v25 = 25;
        int v26 = 26;
        int v27 = 27;
        int v28 = 28;
        int v29 = 29;
        int v30 = 30;
        int v31 = 31;
        int v32 = 32;
        int v33 = 33;
        int v34 = 34;
        int v35 = 35;
        int v36 = 36;
        int v37 = 37;
        int v38 = 38;
        int v39 = 39;
        int v40 = 40;
        int v41 = 41;
        int v42 = 42;
        int v43 = 43;
        int v44 = 44;
        int v45 = 45;
        int v46 = 46;
        int v47 = 47;
        int v48 = 48;
        int v49 = 49;
        int v50 = 50;
        int v51 = 51;
        int v52 = 52;
        int v53 = 53;
        int v54 = 54;
        int v55 = 55;
        int v56 = 56;
        int v57 = 57;
        int v58 = 58;
        int v59 = 59;
        int v60 = 60;
        int v61 = 61;
        int v62 = 62;
        int v63 = 63;
        int v64 = 64;
        int v65 = 65;
        int v66 = 66;
        int v67 = 67;
        int v68 = 68;
        int v69 = 69;
        int v70 = 70;
        int v71 = 71;
        int v72 = 72;
        int v73 = 73;
        int v74 = 74;
        int v75 = 75;
        int v76 = 76;
        int v77 = 77;
        int v78 = 78;
        int v79 = 79;
        int v80 = 80;
        int v81 = 81;
        int v82 = 82;
        int v83 = 83;
        int v84 = 84;
        int v85 = 85;
        int v86 = 86;
        int v87 = 87;
        int v88 = 88;
        int v89 = 89;
        int v90 = 90;
        int v91 = 91;
        int v92 = 92;
        int v93 = 93;
        int v94 = 94;
        int v95 = 95;
        int v96 = 96;
        int v97 = 97;
        int v98 = 98;
        int v99 = 99;
        int v100 = 100;
        int v101 = 101;
        int v102 = 102;
        int v103 = 103;
        int v104 = 104;
        int v105 = 105;
        int v106 = 106;
        int v107 = 107;
        int v108 = 108;
        int v109 = 109;
        int v110 = 110;
        int v111 = 111;
        int v112 = 112;
        int v113 = 113;
        int v114 = 114;
        int v115 = 115;
        int v116 = 116;
        int v117 = 117;
        int v118 = 118;
        int v119 = 119;
        int v120 = 120;
        int v121 = 121;
        int v122 = 122;
        int v123 = 123;
        int v124 = 124;
        int v125 = 125;
        int v126 = 126;
        int v127 = 127;
        int v128 = 128;
        int v129 = 129;
        int v130 = 130;
        int v131 = 131;
        int v132 = 132;
        int v133 = 133;
        int v134 = 134;
        int v135 = 135;
        int v136 = 136;
        int v137 = 137;
        int v138 = 138;
        int v139 = 139;
        int v140 = 140;
        int v141 = 141;
        int v142 = 142;
        int v143 = 143;
        int v144 = 144;
        int v145 = 145;
        int v146 = 146;
        int v147 = 147;
        int v148 = 148;
        int v149 = 149;
        int v150 = 150;
        int v151 = 151;
        int v152 = 152;
        int v153 = 153;
        int v154 = 154;
        int v155 = 155;
        int v156 = 156;
        int v157 = 157;
        int v158 = 158;
        int v159 = 159;
        int v160 = 160;
        int v161 = 161;
        int v162 = 162;
        int v163 = 163;
        int v164 = 164;
        int v165 = 165;
        int v166 = 166;
        int v167 = 167;
        int v168 = 168;
        int v169 = 169;
        int v170 = 170;
        int v171 = 171;
        int v172 = 172;
        int v173 = 173;
        int v174 = 174;
        int v175 = 175;
        int v176 = 176;
        int v177 = 177;
        int v178 = 178;
        int v179 = 179;
        int v180 = 180;
        int v181 = 181;
        int v182 = 182;
        int v183 = 183;
        int v184 = 184;
        int v185 = 185;
        int v186 = 186;
        int v187 = 187;
        int v188 = 188;
        int v189 = 189;
        int v190 = 190;
        int v191 = 191;
        int v192 = 192;
        int v193 = 193;
        int v194 = 194;
        int v195 = 195;
        int v196 = 196;
        int v197 = 197;
        int v198 = 198;
        int v199 = 199;
        int v200 = 200;
        int v201 = 201;
        int v202 = 202;
        int v203 = 203;
        int v204 = 204;
        int v205 = 205;
        int v206 = 206;
        int v207 = 207;
        int v208 = 208;
        int v209 = 209;
        int v210 = 210;
        int v211 = 211;
        int v212 = 212;
        int v213 = 213;
        int v214 = 214;
        int v215 = 215;
        int v216 = 216;
        int v217 = 217;
        int v218 = 218;
        int v219 = 219;
        int v220 = 220;
        int v221 = 221;
        int v222 = 222;
        int v223 = 223;
        int v224 = 224;
        int v225 = 225;
        int v226 = 226;
        int v227 = 227;
        int v228 = 228;
        int v229 = 229;
        int v230 = 230;
        int v231 = 231;
        int v232 = 232;
        int v233 = 233;
        int v234 = 234;
        int v235 = 235;
        int v236 = 236;
        int v237 = 237;
        int v238 = 238;
        int v239 = 239;
        int v240 = 240;
        int v241 = 241;
        int v242 = 242;
        int v243 = 243;
        int v244 = 244;
        int v245 = 245;
        int v246 = 246;
        int v247 = 247;
        int v248 = 248;
        int v249 = 249;
        int v250 = 250;
        int v251 = 251;
        int v252 = 252;
        int v253 = 253;
        int v254 = 254;
        int v255 = 255;
        int v256 = 256;
        int v257 = 257;
        int v258 = 258;
        int v259 = 259;
        int v260 = 260;
        int v261 = 261;
        int v262 = 262;
        int v263 = 263;
        int v264 = 264;
        int v265 = 265;
        int v266 = 266;
        int v267 = 267;
        int v268 = 268;
        int v269 = 269;
        int v270 = 270;
        int v271 = 271;
        int v272 = 272;
        int v273 = 273;
        int v274 = 274;
        int v275 = 275;
        int v276 = 276;
        int v277 = 277;
        int v278 = 278;
        int v279 = 279;
        int v280 = 280;
        int v281 = 281;
        int v282 = 282;
        int v283 = 283;
        int v284 = 284;
        int v285 = 285;
        int v286 = 286;
        int v287 = 287;
        int v288 = 288;
        int v289 = 289;
        int v290 = 290;
        int v291 = 291;
        int v292 = 292;
        int v293 = 293;
        int v294 = 294;
        int v295 = 295;
        int v296 = 296;
        int v297 = 297;
        int v298 = 298;
        int v299 = 299;
        int acc = seed + v0 + v1;
        v268 += 5;
        acc += v268;
        acc += v260;
        v298 = v298 + 1000;
        acc += v298;
        acc += v299;
        System.out.println(acc);
        return acc;
    }

    static long wideLongs(long seed) {
        long x0 = 5000L;
        long x1 = 5001L;
        long x2 = 5002L;
        long x3 = 5003L;
        long x4 = 5004L;
        long x5 = 5005L;
        long x6 = 5006L;
        long x7 = 5007L;
        long x8 = 5008L;
        long x9 = 5009L;
        long x10 = 5010L;
        long x11 = 5011L;
        long x12 = 5012L;
        long x13 = 5013L;
        long x14 = 5014L;
        long x15 = 5015L;
        long x16 = 5016L;
        long x17 = 5017L;
        long x18 = 5018L;
        long x19 = 5019L;
        long x20 = 5020L;
        long x21 = 5021L;
        long x22 = 5022L;
        long x23 = 5023L;
        long x24 = 5024L;
        long x25 = 5025L;
        long x26 = 5026L;
        long x27 = 5027L;
        long x28 = 5028L;
        long x29 = 5029L;
        long x30 = 5030L;
        long x31 = 5031L;
        long x32 = 5032L;
        long x33 = 5033L;
        long x34 = 5034L;
        long x35 = 5035L;
        long x36 = 5036L;
        long x37 = 5037L;
        long x38 = 5038L;
        long x39 = 5039L;
        long x40 = 5040L;
        long x41 = 5041L;
        long x42 = 5042L;
        long x43 = 5043L;
        long x44 = 5044L;
        long x45 = 5045L;
        long x46 = 5046L;
        long x47 = 5047L;
        long x48 = 5048L;
        long x49 = 5049L;
        long x50 = 5050L;
        long x51 = 5051L;
        long x52 = 5052L;
        long x53 = 5053L;
        long x54 = 5054L;
        long x55 = 5055L;
        long x56 = 5056L;
        long x57 = 5057L;
        long x58 = 5058L;
        long x59 = 5059L;
        long x60 = 5060L;
        long x61 = 5061L;
        long x62 = 5062L;
        long x63 = 5063L;
        long x64 = 5064L;
        long x65 = 5065L;
        long x66 = 5066L;
        long x67 = 5067L;
        long x68 = 5068L;
        long x69 = 5069L;
        long x70 = 5070L;
        long x71 = 5071L;
        long x72 = 5072L;
        long x73 = 5073L;
        long x74 = 5074L;
        long x75 = 5075L;
        long x76 = 5076L;
        long x77 = 5077L;
        long x78 = 5078L;
        long x79 = 5079L;
        long x80 = 5080L;
        long x81 = 5081L;
        long x82 = 5082L;
        long x83 = 5083L;
        long x84 = 5084L;
        long x85 = 5085L;
        long x86 = 5086L;
        long x87 = 5087L;
        long x88 = 5088L;
        long x89 = 5089L;
        long x90 = 5090L;
        long x91 = 5091L;
        long x92 = 5092L;
        long x93 = 5093L;
        long x94 = 5094L;
        long x95 = 5095L;
        long x96 = 5096L;
        long x97 = 5097L;
        long x98 = 5098L;
        long x99 = 5099L;
        long x100 = 5100L;
        long x101 = 5101L;
        long x102 = 5102L;
        long x103 = 5103L;
        long x104 = 5104L;
        long x105 = 5105L;
        long x106 = 5106L;
        long x107 = 5107L;
        long x108 = 5108L;
        long x109 = 5109L;
        long x110 = 5110L;
        long x111 = 5111L;
        long x112 = 5112L;
        long x113 = 5113L;
        long x114 = 5114L;
        long x115 = 5115L;
        long x116 = 5116L;
        long x117 = 5117L;
        long x118 = 5118L;
        long x119 = 5119L;
        long x120 = 5120L;
        long x121 = 5121L;
        long x122 = 5122L;
        long x123 = 5123L;
        long x124 = 5124L;
        long x125 = 5125L;
        long x126 = 5126L;
        long x127 = 5127L;
        long x128 = 5128L;
        long x129 = 5129L;
        long acc = seed + x128 + x129;
        x127 = x127 + 7L;
        acc += x127;
        x128 = x128 + 11L;
        acc += x128;
        x129 = x129 + 13L;
        acc += x129;
        System.out.println(acc);
        return acc;
    }

    static int neverSwitch(int key) {
        switch (key) {
            case 1: return 1;
            case 2: return 2;
            default: return 0;
        }
    }

    public static void main(String[] args) {
        System.out.println(dense(3));
        System.out.println(dense(7));
        System.out.println(dense(11));
        System.out.println(offsetDense(15));
        System.out.println(offsetDense(0));
        System.out.println(sparse(1000000));
        System.out.println(nested(0, 2));
        System.out.println(loopSwitch(9));
        System.out.println(charSwitch('a'));
        System.out.println(charSwitch('q'));
        System.out.println(byteSwitch((byte) -2));
        System.out.println(byteSwitch((byte) 9));
        System.out.println(shortSwitch((short) -300));
        System.out.println(shortSwitch((short) 301));
        System.out.println(fallthrough(5));
        System.out.println(fallthrough(4));
        System.out.println(fallthrough(7));
        System.out.println(fallthrough(9));
        System.out.println(noDefault(0));
        System.out.println(noDefault(1));
        System.out.println(noDefault(2));
        System.out.println(deepIf(20000));
        System.out.println(deepIf(700));
        System.out.println(deepIf(20));
        System.out.println(deepIf(5));
        System.out.println(deepIf(0));
        System.out.println(deepIf(-3));
        System.out.println(wideInts(7));
        System.out.println(wideLongs(9L));
    }
}
