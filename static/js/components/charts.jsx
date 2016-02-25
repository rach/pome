import d3 from 'd3';
import React from 'react';
import ReactFauxDOM from 'react-faux-dom';
import memoize from 'memoize-decorator';
import JQuery from 'jquery';
import {formatPercent, formatBytes} from '../utils';


export class BarChart extends React.Component {
  static propTypes = {
    data: React.PropTypes.array,
    title: React.PropTypes.string
  }

  margin(){
    return {top: 40, right: 10, bottom: 50, left: 60}
  }

  height(){
    let margin = this.margin();
    return 180 - margin.top - margin.bottom;
  }

  width(){
    let margin = this.margin();
    return  910 - margin.left - margin.right;
  }

  sample(){
    return 120;
  }

  yAxis(){
    return d3.svg.axis()
                  .scale(this.yScale())
                  .ticks(3)
                  .tickSize(2, 0)
                  .tickPadding(6)
                  .tickFormat(this.yFormatter())
                  .orient("left");
  }

  xAxis(){
    return d3.svg.axis()
             .scale(this.xScale())
             .tickValues(d3.range(0, this.props.data.length, 12))
             .tickSize(2, 0)
             .tickPadding(6)
             .tickFormat(this.xFormatter())
             .orient("bottom");
  }

  yScale(){
    var y = d3.scale.linear()
              .domain([0, this.yMax()])
              .range([this.height(), 0]);
    return y;
  }

  xScale(){
    var x = d3.scale.ordinal()
              .domain(d3.range(this.sample()))
              .rangeRoundBands([this.width(), 0], 0.20, 0.10);
    return x;
  }

  yMax(){
      return d3.max(this.props.data, v => v.y) || 10;
  }

  xFormatter(){
    const timeFormat = d3.time.format("%I:%M %p");
    let formatter = t => {
      timeFormat(new Date(this.props.data[t].x * 1000));
    };
    return formatter;
  }

  yFormatter(){
    return v => v;
  }

  render() {
    var that = this;
    var data = this.props.data,
        xdata = this.props.x;

    var y = this.yScale();
    var x = this.xScale();

    let detail = (<span></span>);
    const datetimeFormat = d3.time.format("%d/%m/%y %I:%M %p");

    const xFormatter = this.xFormatter();
    const yFormatter = this.yFormatter();

    const xAxis = this.xAxis();
    const yAxis = this.yAxis();
    const node = ReactFauxDOM.createElement('svg');
    const margin = this.margin();
    const width = this.width();
    const height  = this.height();
    const svg = d3.select(node)
                  .attr("width", width + margin.left + margin.right)
                  .attr("height", height + margin.top + margin.bottom)
                  .append("g")
                  .attr("transform", "translate(0," + margin.top + ")");

    var rect = svg.selectAll("rect")
                  .data(data)
                  .enter().append("rect")
                  .attr("x", function(d, i) { return x(i) + margin.left; })
                  .attr("width", x.rangeBand())
                  .attr("y", function(d) { return y(d.y); })
                  .attr("height", function(d) { return y(0) - y(d.y); })
                  .on("mouseover", (d, i) => {
                    //this should use action dispatch to update the state
                    var suffix = "";
                    if (d.context){
                      suffix = "<br/> " + d.context;
                    }

                    var t = datetimeFormat(new Date(d.x * 1000));
                    JQuery(React.findDOMNode(that))
                                .find('.bar-value')
                                .html(t + " → " + yFormatter(d.y) + suffix);
                  })
                  .on("mouseout", (d, i) => {
                    //this should use action dispatch to update the state
                    JQuery(React.findDOMNode(that)).find('.bar-value').text("");
                  });

    svg.append("g")
       .attr("class", "x axis")
       .attr("transform", "translate("+ margin.left+"," + height + ")")
       .call(xAxis)
       .selectAll("text")
       .attr("dx", "-.8em")
       .attr("dy", ".15em")
       .style("text-anchor", "end")
       .attr("transform", "rotate(-65)" );

    svg.append("g")
       .attr("class", "y axis")
       .attr("transform", "translate("+ margin.left + ", 0)")
       .call(yAxis);

    var subtitle = "";
    if (this.props.subtitle){
      subtitle = (<h6>{this.props.subtitle}</h6>);

    }
    return (
      <div className="chart">
        <div className="bar-value text-right">
        </div>
        <div className="row">
          <div className="col-sm-8">
            <h5>{this.props.title}</h5>
          </div>
        </div>
        <div className="row">
          <div className="col-sm-12">
            {subtitle}
          </div>
        </div>
        {node.toReact()}
      </div>
    );

  }
}

function roundBytes(bytes){
  if(bytes == 0) return '0 Byte';
  var k = 1000;
  var i = Math.floor(Math.log(bytes) / Math.log(k));
  return Math.ceil(bytes / Math.pow(k, i)) * Math.pow(k, i);
}

export class BytesBarChart extends BarChart {
  yFormatter(){
    return formatBytes;
  }
  yMax(){
    return d3.max(this.props.data, v => v.y) || 1000000;
  }
  yAxis(){
    const max_bytes = this.yMax();
    const max_thick_val = roundBytes(max_bytes);
    const middle_thick_val = roundBytes(Math.floor(max_bytes / 2));

    return d3.svg.axis()
             .scale(this.yScale())
             .ticks(3)
             .tickValues([0, middle_thick_val, max_thick_val])
             .tickSize(2, 0)
             .tickPadding(6)
             .tickFormat(this.yFormatter())
             .orient("left");
  }
  yScale(){
    const max_bytes = this.yMax();
    const max_thick_val = roundBytes(max_bytes);
    var y = d3.scale.linear()
              .domain([0, max_thick_val])
              .range([this.height(), 0]);
    return y;
  }
}

export class PercentBarChart extends BarChart {
  yFormatter(){
    return formatPercent;
  }
  yMax(){
    return 100;
  }
}

