import { Component } from '@angular/core';
import { PagingController } from '../common/paging';
import { HttpClient } from '@angular/common/http';
import { Toaster } from '../toaster.service';

export interface BuildInfo {
  id?: number
  name?: string
  status?: string
  ctime?: number | Date
  utime?: number | Date
  buildSteps?: string[]
}

@Component({
  selector: 'app-build-info',
  templateUrl: './build-info.component.html',
  styleUrl: './build-info.component.css'
})
export class BuildInfoComponent {

  data: any[] = []
  pagingController: PagingController;

  constructor(private http: HttpClient, private toaster: Toaster) {
  }

  fetchList() {
    this.http.post<any>("/api/build/info/list", this.pagingController.paging)
      .subscribe({
        next: (resp) => {
          if (resp.error){
            this.toaster.toast(resp.msg);
            return;
          }

          this.data = [];
          if (resp.data.payload) {
            for (let r of resp.data.payload) {
              if (r.ctime) {
                r.ctime = new Date(r.ctime);
              }
              if (r.utime) {
                r.utime = new Date(r.utime);
              }
              this.data.push(r);
            }
          }
          this.pagingController.onTotalChanged(resp.data.paging);
        },
        error: (err) => {
          this.toaster.toast(`Request failed: ${err}`)
        },
      });
  }

  onPagingControllerReady(pc: PagingController) {
    this.pagingController = pc;
    this.pagingController.onPageChanged = () => this.fetchList();
    this.fetchList();
  }

  triggerBuild(u : BuildInfo) {
    this.http.post<any>("/api/build/trigger", { name : u.name })
      .subscribe({
        next: (resp) => {
          if (resp.error){
            this.toaster.toast(resp.msg);
            return;
          }
          this.fetchList();
        },
        error: (err) => {
          this.toaster.toast(`Request failed: ${err}`)
        },
      });
  }
}
