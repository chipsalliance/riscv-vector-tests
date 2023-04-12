/* Ref: https://github.com/OpenRigil/openrigil-firmware/blob/master/bootrom.ld
 */

OUTPUT_ARCH( "riscv" )
ENTRY(_start)

MEMORY
{
    BOOTROM (rx) : ORIGIN = 0x100000, LENGTH = 288K
    RAM    (xrw) : ORIGIN = 0x80000000, LENGTH = 64K
}

SECTIONS
{
  . = 0x1000;
  .text.init : {
    *(.text.init)
  } >BOOTROM

  . = ALIGN(0x1000);
  .text : {
    *(.text)
  } >BOOTROM

  . = ALIGN(0x1000);
  .data : {
    *(.data)
    PROVIDE (__data_end = .);
  } >RAM AT> BOOTROM

  PROVIDE (__data_start = ADDR(.data));
  PROVIDE (__data_load_start = LOADADDR(.data));

  _end = .;
}
