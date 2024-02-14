// This main runs RISC-V binaries in spike, but generates test cases
// on executing a custom0 instruction

#include <riscv/sim.h>
#include <riscv/extension.h>
#include <fesvr/option_parser.h>

// Copied from spike main.
// TODO: This should really be provided in libriscv
static std::vector<std::pair<reg_t, abstract_mem_t*>> make_mems(const std::vector<mem_cfg_t> &layout)
{
  std::vector<std::pair<reg_t, abstract_mem_t*>> mems;
  mems.reserve(layout.size());
  for (const auto &cfg : layout) {
    mems.push_back(std::make_pair(cfg.get_base(), new mem_t(cfg.get_size())));
  }
  return mems;
}

static reg_t magic_insn(processor_t* p, insn_t insn, reg_t pc) {
  static int ncase = 2;
  int group = insn.rs1();
  bool vxsat = insn.rs2() & 0x1;
  for (int reg = group; reg < 2*group; reg++) {
    for (int i = 0; i < p->VU.VLEN / p->get_xlen(); i++) {
      if (p->get_xlen() == 64) {
        printf(
               "  TEST_CASE(%d, t0, 0x%lx, ld t0, 0(a0); addi a0, a0, 8)\n",
               ncase++,
               p->VU.elt<type_sew_t<64>::type>(reg, i, false));
      } else {
        printf(
               "  TEST_CASE(%d, t0, 0x%x, lw t0, 0(a0); addi a0, a0, 4)\n",
               ncase++,
               p->VU.elt<type_sew_t<32>::type>(reg, i, false));
      }
    }
  }
  if (vxsat) {
    printf("  TEST_CASE(%d, t0, 0x%lx, csrr t0, vxsat)\n", ncase++, p->get_csr(0x009));
  }
  printf("---\n");
  return pc + 4;
}

class magic_extension_t : public extension_t {
  std::vector<insn_desc_t> get_instructions() override {
    std::vector<insn_desc_t> insns;
    insns.push_back((insn_desc_t){0x0000000B, 0x0000007F,
                                  magic_insn, magic_insn, magic_insn, magic_insn,
                                  magic_insn, magic_insn, magic_insn, magic_insn});
    return insns;
  }
  std::vector<disasm_insn_t*> get_disasms() override { return std::vector<disasm_insn_t*>(); }
  const char* name() override { return "magic"; }
};

int main(int argc, char** argv) {
  std::vector<mem_cfg_t> mem_cfg { mem_cfg_t(0x80000000, 0x10000000) };
  std::vector<size_t> hartids = {0};
  cfg_t cfg;
  option_parser_t parser;
  parser.option(0, "isa", 1, [&](const char* s){cfg.isa = s;});
  parser.option(0, "varch", 1, [&](const char* s){cfg.varch = s;});

  auto argv1 = parser.parse(argv);
  std::vector<std::string> htif_args(argv1, (const char*const*)argv + argc);

  debug_module_config_t dm_config = {
    .progbufsize = 2,
    .max_sba_data_width = 0,
    .require_authentication = false,
    .abstract_rti = 0,
    .support_hasel = true,
    .support_abstract_csr_access = true,
    .support_abstract_fpr_access = true,
    .support_haltgroups = true,
    .support_impebreak = true
  };
  std::vector<std::pair<reg_t, abstract_mem_t*>> mems = make_mems(cfg.mem_layout);
  std::vector<device_factory_t*> plugin_devices;
  sim_t sim(&cfg, false,
            mems,
            plugin_devices,
            htif_args,
            dm_config,
            nullptr,
            true,
            nullptr,
            false,
            nullptr);
  magic_extension_t magic;
  sim.get_core(0)->register_extension(&magic);
  sim.run();
}
