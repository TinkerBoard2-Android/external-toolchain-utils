import job_manager
from utils import utils
import SimpleXMLRPCServer

class Server:

  def __init__(self):
    self.job_manager = job_manager.JobManager()
    print "Started server thread."


  def ExecuteJobGroup(self, job_group):
    #TODO(raymes): Verify that the job graph is valid. I.e. every
    # dependency has been transmitted to the server.
    job_group = utils.Deserialize(job_group)
    for current_job in job_group:
      self.job_manager.AddJob(current_job)


  def GetAllJobs(self):
    jobs_dict = {}
    jobs_dict["all"] = self.job_manager.all_jobs
    jobs_dict["ready"] = self.job_manager.ready_jobs
    jobs_dict["pending"] = self.job_manager.pending_jobs
    jobs_dict["executing"] = self.job_manager.executing_jobs
    jobs_dict["completed"] = self.job_manager.completed_jobs


  def StartServer(self):
    self.job_manager.StartJobManager()

  def StopServer(self):
    self.job_manager.StopJobManager()

if __name__ == "__main__":
  server = Server()
  server.StartServer()
  xmlserver = SimpleXMLRPCServer.SimpleXMLRPCServer(("localhost", 8000),
                                                  allow_none=True)
  xmlserver.register_instance(server)
  try:
    xmlserver.serve_forever()
  except (KeyboardInterrupt, SystemExit):
    print "Caught exception... Cleaning up."
    server.StopServer()
    raise
    

