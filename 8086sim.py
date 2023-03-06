import sys
from pathlib import Path

REG_ENCODING = {
    "0": {
        "000": "al",
        "001": "cl",
        "010": "dl",
        "011": "bl",
        "100": "ah",
        "101": "ch",
        "110": "dh",
        "111": "bh",
    }, 
    "1": {
        "000": "ax",
        "001": "cx",
        "010": "dx",
        "011": "bx",
        "100": "sp",
        "101": "bp",
        "110": "si",
        "111": "di",
    },
}

RM_ENCODING = {
    "000": "bx + si",
    "001": "bx + di",
    "010": "bp + si",
    "011": "bp + di",
    "100": "si",
    "101": "di",
    "110": "bp", # Note (Tom) watch for direct address
    "111": "bx",
}

MOD_ENCODING = {
    "00": 0, # Note (Tom) watch for rm 110
    "01": 8,
    "10": 16,
    "11": 0,
}

def _calc_dec(ins: str) -> int:
    lsb = ins[0:8]
    num_bytes = 1
    if len(ins) == 16:
        msb = ins[8:16]
        dec = int(msb + lsb, 2)
        num_bytes = 2
    else:
        dec = int(lsb, 2)

    b = dec.to_bytes(num_bytes, byteorder=sys.byteorder, signed=False)
    return int.from_bytes(b, byteorder=sys.byteorder, signed=True)

class Mov:
    @staticmethod
    def reg_mem_to_from_reg(ins: str, stream_index: int) -> tuple[str, int]:
        
        total_bits_read = 16
        
        d = ins[stream_index+6]
        w = ins[stream_index+7]
        
        mod = ins[stream_index+8:stream_index+10]
        
        reg = ins[stream_index+10:stream_index+13]
        rm = ins[stream_index+13:stream_index+16]

        # get displacement
        disp_bits = MOD_ENCODING[mod]
        direct_address = False

        # account for edge case
        if mod == "00" and rm == "110":
            disp_bits = 16
            direct_address = True
            

        # register mode 
        if mod == "11" and disp_bits == 0:
        
            if d == "0":
                source = REG_ENCODING[w][reg]
                dest = REG_ENCODING[w][rm]
            else:
                source = REG_ENCODING[w][rm]
                dest = REG_ENCODING[w][reg]

        # effective address calc
        else:
            # convert displacement to decimal
            dec = 0
            sign = "+"

            if disp_bits > 0:
                dec = _calc_dec(ins[stream_index+16:stream_index+16+disp_bits])
                if dec < 0:
                    sign = "-"
                    dec = abs(dec)
                total_bits_read += disp_bits

            if d == "0":
                source = REG_ENCODING[w][reg]
                if direct_address:
                    dest = f"[{dec}]"
                elif dec != 0:
                    dest = f"[{RM_ENCODING[rm]} {sign} {dec}]"
                else:
                    dest = f"[{RM_ENCODING[rm]}]"
            else:
                if direct_address:
                    source = f"[{dec}]"
                elif dec != 0:
                    source = f"[{RM_ENCODING[rm]} {sign} {dec}]"
                else:
                    source = f"[{RM_ENCODING[rm]}]"
                dest = REG_ENCODING[w][reg]

        return (f"{dest}, {source}"), stream_index + total_bits_read
    
    @staticmethod
    def immediate_to_reg(ins: str, stream_index: int) -> tuple[str, int]:

        total_bits_read = 8

        w = ins[stream_index+4]
        reg = ins[stream_index+5:stream_index+8]

        if w == "1":
            dec = _calc_dec(ins[stream_index+8:stream_index+24])
            total_bits_read += 16
        else:
            dec = _calc_dec(ins[stream_index+8:stream_index+16])
            total_bits_read += 8
        
        return (f"{REG_ENCODING[w][reg]}, {dec}", stream_index + total_bits_read)    
    
    @staticmethod
    def mem_acc(ins: str, stream_index: int, to_ax: bool = True) -> tuple[str, int]:

        total_bits_read = 16

        w = ins[stream_index+7]

        if w == "1":
            dec = _calc_dec(ins[stream_index+8:stream_index+24])
            total_bits_read += 8
        else:
            dec = _calc_dec(ins[stream_index+8:stream_index+16],)
        
        if to_ax:
            return (f"ax, [{dec}]"), stream_index + total_bits_read
        else:
            return (f"[{dec}], ax"), stream_index + total_bits_read
    
    @staticmethod
    def immediate_to_reg_mem(ins: str, stream_index: int) -> tuple[str, int]:
        
        total_bits_read = 16

        w = ins[stream_index+7]
        mod = ins[stream_index+8:stream_index+10]
        rm = ins[stream_index+13:stream_index+16]

        # get displacement
        disp_bits = MOD_ENCODING[mod]
        direct_address = False

        # account for edge case
        if mod == "00" and rm == "110":
            disp_bits = 16
            direct_address = True
        
        dec = 0
        if disp_bits > 0:
            dec = _calc_dec(ins[stream_index+16:stream_index+16+disp_bits])
            total_bits_read += disp_bits

        dest = f"[{RM_ENCODING[rm]}]"
        if direct_address:
            dest = f"[{dec}]"
        elif dec != 0:
            dest = f"[{RM_ENCODING[rm]} + {dec}]"

        if w == "1":
            val = f"word {_calc_dec(ins[stream_index+total_bits_read:stream_index+total_bits_read+16])}"
            total_bits_read += 16
        else:
            val = f"byte {_calc_dec(ins[stream_index+total_bits_read:stream_index+total_bits_read+8])}"
            total_bits_read += 8

        return (f"{dest}, {val}"), stream_index + total_bits_read


def decode(stream: str) -> None:
    """
    decode the instruction stream and print the results
    """

    stream_index = 0

    while stream_index < len(stream) - 1:
        ## MOV
        # check the first 4 bits of the instruction stream 
        if stream[stream_index: stream_index+4] == "1011":
            # immediate to register
            decoded, stream_index = Mov.immediate_to_reg(stream, stream_index)
            print(f"mov {decoded}")
        elif stream[stream_index: stream_index+6] == "100010":
            # register/memory to/from register
            decoded, stream_index = Mov.reg_mem_to_from_reg(stream, stream_index)
            print(f"mov {decoded}")
        elif stream[stream_index: stream_index+7] == "1100011":
            # immediate to register/memory
            decoded, stream_index = Mov.immediate_to_reg_mem(stream, stream_index)
            print(f"mov {decoded}")
        elif stream[stream_index: stream_index+7] == "1010001":
            # accumulator to memory
            decoded, stream_index = Mov.mem_acc(stream, stream_index, to_ax=False)
            print(f"mov {decoded}")
        elif stream[stream_index: stream_index+7] == "1010000":
            # memory to accumulator
            decoded, stream_index = Mov.mem_acc(stream, stream_index)
            print(f"mov {decoded}")
            

if __name__ == "__main__":
    # preprocess binary file to get instruction stream
    fn = sys.argv[1]

    raw = Path(fn).read_bytes()
    
    stream = "".join(format(x, '08b') for x in bytearray(raw))
    assert len(stream) % 8 == 0 

    decode(stream)