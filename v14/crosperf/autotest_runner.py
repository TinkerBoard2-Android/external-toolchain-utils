#!/usr/bin/python
#
# Copyright 2011 Google Inc. All Rights Reserved.

from utils import command_executer


class AutotestRunner(object):
  """ This defines the interface from crosperf to ./run_remote_tests.sh.
  """
  def __init__(self, logger_to_use=None):
    self._logger = logger_to_use
    self._ce = command_executer.GetCommandExecuter(self._logger)
    self._ct = command_executer.CommandTerminator()

  def Run(self, machine_name, chromeos_root, board, autotest_name,
          autotest_args):
    options = ""
    if board:
      options += " --board=%s" % board
    if autotest_args:
      options += " %s" % autotest_args
    command = ("./run_remote_tests.sh --remote=%s %s %s" %
               (machine_name, options, autotest_name))
    return self._ce.ChrootRunCommand(chromeos_root, command, True, self._ct)

  def Terminate(self):
    self._ct.Terminate()


class MockAutotestRunner(object):
  def __init__(self):
    pass

  def Run(self, *args):
    return ["", "", 0]
