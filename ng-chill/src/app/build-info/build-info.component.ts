import { Component } from '@angular/core';
import { PagingController } from '../common/paging';
import { HttpClient } from '@angular/common/http';
import { Toaster } from '../toaster.service';
import { NavigationService } from '../navigation.service';

export interface ApiListBuildInfoRes {
  id?: number                    // build info id
  name?: string                  // build name
  status?: string                // last build status
  ctime?: number                 // create time
  utime?: number                 // update time
  commitId?: string              // last build commit id
  tag?: string                   // tag
  buildSteps?: string[]          // build steps
  triggerable?: boolean          // whether the build is triggerable
}

@Component({
  selector: 'app-build-info',
  templateUrl: './build-info.component.html',
  styleUrl: './build-info.component.css'
})
export class BuildInfoComponent {

  data: ApiListBuildInfoRes[] = []
  pagingController: PagingController;

  constructor(private http: HttpClient, private toaster: Toaster, private nav: NavigationService) {
  }

  fetchList() {
    this.http.post<any>("/api/build/info/list", this.pagingController.paging)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.toaster.toast(resp.msg);
            return;
          }

          this.data = [];
          if (resp.data.payload) {
            this.data = resp.data.payload;
          }
          this.pagingController.onTotalChanged(resp.data.paging);
        },
        error: (err) => {
          this.toaster.toast(`Request failed`)
          console.log(err)
        },
      });
  }

  onPagingControllerReady(pc: PagingController) {
    this.pagingController = pc;
    this.pagingController.onPageChanged = () => this.fetchList();
    this.fetchList();
  }

  triggerBuild(u: ApiListBuildInfoRes) {
    this.http.post<any>("/api/build/trigger", { name: u.name })
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.toaster.toast(resp.msg);
            return;
          }
          this.fetchList();
        },
        error: (err) => {
          this.toaster.toast(`Request failed`)
          console.log(err)
        },
      });
  }

  redirectBuildHistory(u: ApiListBuildInfoRes) {
    this.nav.navigateToUrl("/build/history/list", [
      { name: u.name },
    ]);
  }
}
