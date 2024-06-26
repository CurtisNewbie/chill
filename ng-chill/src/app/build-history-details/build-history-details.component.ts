import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { Toaster } from '../toaster.service';
import { ActivatedRoute } from '@angular/router';

export interface ApiCmdLogRes {
  id?: number
  command?: string
  remark?: string
  status?: string
}

export interface ApiQryBuildHistDetailRes {
  id?: number                    // build history id
  name?: string                  // build name
  buildNo?: string               // build no
  commitId?: string              // build commit id
  tag?: string                   // tag
  status?: string                // built status
  startTime?: number             // build start time
  endTime?: number               // build end time
  remark?: string                // remark
  commandLogs?: ApiCmdLogRes[]
}

@Component({
  selector: 'app-build-history-details',
  templateUrl: './build-history-details.component.html',
  styleUrl: './build-history-details.component.css'
})
export class BuildHistoryDetailsComponent implements OnInit {

  data: ApiQryBuildHistDetailRes = {};
  buildNo: string;

  constructor(private http: HttpClient, private toaster: Toaster, private route: ActivatedRoute) {
  }

  ngOnInit(): void {
    this.route.paramMap.subscribe((params) => {
      this.buildNo = params.get("buildNo");
      this.fetchDetail();
    });
  }

  fetchDetail() {
    this.http.post<any>("/api/build/history/detail", { buildNo: this.buildNo })
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.toaster.toast(resp.msg);
            return;
          }

          this.data = resp.data;
        },
        error: (err) => {
          this.toaster.toast(`Request failed`)
          console.log(err)
        },
      });
  }
}
